#version 460 core

#define PI 3.141592653589793 
#define MAX_BOUNCE_COUNT 8 
#define NUM_RAYS_PER_PIXEL 4

const float INF = 1e30;
const int STACK_SIZE = 64;

layout(local_size_x = 8, local_size_y = 8) in;

layout(rgba32f, binding = 0) uniform image2D renderImage;
layout(rgba32f, binding = 1) uniform image2D accumulation;

uniform mat4 invViewMat;
uniform mat4 invProjMat;
uniform vec3 camPos;

uniform int iFrame;
uniform int sphereCount;
uniform int triangleCount;

uniform vec3 sunDirection;
uniform vec3 sunColor;
uniform float sunIntensity;

struct Sphere{
  // XYZ -> Center, W -> Radius
  vec4 pos;
  // X -> Material Index
  vec4 material;
};

struct Triangle{
  vec4 posA, posB, posC;
  vec4 normalA, normalB, normalC;
  // X -> Material Index
  vec4 material;
};

struct RayTracingMaterial{
  vec4 colour;
  // XYZ -> Colour, W -> Strength
  vec4 emission;
};

struct Node {
  vec4 minBound;
  vec4 maxBound;
	// nodeInfo :  X -> Left, Y -> Right, Z -> First Primitive, W -> Primitive Count
  vec4 nodeInfo;
};

struct Primitive { 
	// ssboInfo : X -> Type, Y -> SSBO Index
  vec4 ssboInfo;
};

layout(std430, binding = 2) buffer Spheres {
  Sphere spheres[];
};
layout(std430, binding = 3) buffer Materials{
  RayTracingMaterial materials[];
};
layout(std430, binding = 4) buffer Triangles{
  Triangle triangles[];
};
layout(std430, binding = 5) buffer Nodes {
  Node nodes[];
};
layout(std430, binding = 6) buffer Primitives{
  Primitive primitives[];
};
layout(std430, binding = 7) buffer PrimitivesIndices{
  int primitiveIndices[];
};


struct Ray{
  vec3 origin;
  vec3 dir;
};

struct HitInfo {
  int didHit;
  float dst;
  vec3 hitPoint;
  vec3 normal;
  RayTracingMaterial material;
};

int wang_hash(inout int seed)
{
    seed = int(seed ^ int(61)) ^ int(seed >> int(16));
    seed *= int(9);
    seed = seed ^ (seed >> 4);
    seed *= int(0x27d4eb2d);
    seed = seed ^ (seed >> 15);
    return seed;
}

float RandomFloat01(inout int state)
{
    return float(wang_hash(state)) / 4294967296.0;
}

vec3 RandomUnitVector(inout int state)
{
    float z = RandomFloat01(state) * 2.0f - 1.0f;
    float a = RandomFloat01(state) * 2.0 * PI;
    float r = sqrt(1.0f - z * z);
    float x = r * cos(a);
    float y = r * sin(a);
    return vec3(x, y, z);
}

vec3 Sky(Ray ray){
    float skyGradientT = pow(smoothstep(0.0, 0.4, ray.dir.y), 0.35);
    float groundToSkyT = smoothstep(-0.01, 0.0, ray.dir.y);

    vec3 horizonColour = vec3(0.9, 0.95, 1.0);
    vec3 zenithColour = vec3(0.3, 0.5, 1.0);
    vec3 groundColour = vec3(0.25, 0.25, 0.25);

    vec3 skyGradient = mix(horizonColour, zenithColour, skyGradientT);

    vec3 sunDir = normalize(vec3(10.0, 20.0, 5.0));
    float sunFocus = 100.0;
    float sunIntensity = 200.0;

    float sun = pow(max(0.0, dot(ray.dir, sunDir)), sunFocus) * sunIntensity;

    vec3 composite = mix(groundColour, skyGradient, groundToSkyT);
    return composite + vec3(sun);
}

void RaySphere(Ray ray, Sphere sphere, inout HitInfo hitInfo){
  vec3 sphereCenter = sphere.pos.xyz;
  float sphereRadius = sphere.pos.w;
  vec3 oc = ray.origin - sphereCenter;

  float a = dot(ray.dir, ray.dir);
  float b = 2.0 * dot(oc, ray.dir);
  float c = dot(oc, oc) - sphereRadius * sphereRadius;

  float discriminant = b * b - 4.0 * a * c;

  if(discriminant >= 0.0){
    float dst = (-b - sqrt(discriminant)) / (2.0 * a);

    if(dst >= 0.0 && dst < hitInfo.dst){
      hitInfo.didHit = 1;
      hitInfo.dst = dst;
      hitInfo.hitPoint = ray.origin + ray.dir * dst;
      hitInfo.normal = normalize(hitInfo.hitPoint - sphereCenter);

      int materialIdx = int(sphere.material.x);
      RayTracingMaterial material = materials[materialIdx];
      hitInfo.material = material;
    }
  }
}

void RayTriangle(Ray ray, Triangle tri, inout HitInfo hitInfo){

  vec3 edgeAB = (tri.posB - tri.posA).xyz;
  vec3 edgeAC = (tri.posC - tri.posA).xyz;
  
  vec3 normalVec = cross(edgeAB, edgeAC);
  vec3 ao = ray.origin - tri.posA.xyz;
  vec3 dao = cross(ao, ray.dir);

  float determinant = -dot(ray.dir, normalVec);
  float invDet = 1.0 / determinant;

  float dst = dot(ao, normalVec) * invDet;
  float u = dot(edgeAC, dao) * invDet;
  float v = -dot(edgeAB, dao) * invDet;
  float w = 1.0 - u - v;

  int didHit = (determinant >= 1E-6 && dst >= 0 && u >= 0 && v >= 0 && w >= 0) ? 1 : 0;

  if (didHit > 0 && dst < hitInfo.dst) {
    hitInfo.didHit = didHit;
    hitInfo.hitPoint = ray.origin + ray.dir * dst;
    hitInfo.normal = normalize(tri.normalA * w + tri.normalB * u + tri.normalC * v).xyz;
    hitInfo.dst = dst;
    
    int materialIdx = int(tri.material.x);
    RayTracingMaterial material = materials[materialIdx];

    hitInfo.material = material;
  }
}

void IntersectPrimitive(Ray ray, Primitive primitive, inout HitInfo closest){

  int type = int(primitive.ssboInfo.x);
  int idx = int(primitive.ssboInfo.y);
  switch (type){
    case 0:
      RaySphere(ray, spheres[idx], closest);
      break;
    case 1:
      RayTriangle(ray, triangles[idx], closest);
      break;
  }
}

bool RayAABB(Ray ray, vec3 minBound, vec3 maxBound, float maxDistance){
  vec3 invDir = 1.0 / ray.dir;
  vec3 t0 = (minBound - ray.origin) * invDir;
  vec3 t1 = (maxBound - ray.origin) * invDir;

  vec3 tmin = min(t0, t1);
  vec3 tmax = max(t0, t1);

  float nearT = max(max(tmin.x, tmin.y), tmin.z);
  float farT  = min(min(tmax.x, tmax.y), tmax.z);

  return ( farT >= max(nearT, 0.0) && nearT < maxDistance);
}


HitInfo TraverseBVH(Ray ray){
  HitInfo closestHit;
  closestHit.didHit = 0;
  closestHit.dst = INF;

  int stack[STACK_SIZE];
  int stackPtr = 0;

  // Root Node
  stack[stackPtr++] = 0;

  while(stackPtr > 0){
    int nodeIdx = stack[--stackPtr];
    Node node = nodes[nodeIdx];

    if(!RayAABB(ray, node.minBound.xyz, node.maxBound.xyz, closestHit.dst)) { 
      continue;
    }

    if(node.nodeInfo.w > 0) {
      int first = int(node.nodeInfo.z);
      int count = int(node.nodeInfo.w);

      for(int i = 0; i < count; i++){
        int primitiveIdx = primitiveIndices[first + i];
        Primitive primitive = primitives[primitiveIdx];
        IntersectPrimitive(ray, primitive, closestHit);
      }
    }
    else{
      stack[stackPtr++] = int(node.nodeInfo.x);
      stack[stackPtr++] = int(node.nodeInfo.y);
    }
  }
  return closestHit;
}

HitInfo CalculateRayCollision(Ray ray){

  HitInfo closestHit;
  closestHit.didHit = 0;
  closestHit.dst = 10000.0;
  RayTracingMaterial initMaterial = {vec4(vec3(0.0), 1.0), vec4(0.0)};
  closestHit.material = initMaterial;

  for(int i = 0; i < int(sphereCount); i++){
    Sphere sphere = spheres[i];
    RaySphere(ray, sphere, closestHit);
  }

  for(int i = 0; i < int(triangleCount); i++){
    Triangle tri = triangles[i];
    RayTriangle(ray, tri, closestHit);
  }
  return closestHit;
}

vec3 TraceRay(Ray ray, inout int rngState){

  vec3 colour = vec3(1.0);
  vec3 incomingLight = vec3(0.0);

  for(int i = 0; i < MAX_BOUNCE_COUNT; i++){
    //HitInfo hitInfo = CalculateRayCollision(ray);
    HitInfo hitInfo = TraverseBVH(ray);

    if(hitInfo.didHit > 0){
      ray.origin = hitInfo.hitPoint;
      ray.dir = normalize(hitInfo.normal + RandomUnitVector(rngState));

      RayTracingMaterial material = hitInfo.material;
      vec3 emittedLight = material.emission.xyz * material.emission.w;
      incomingLight += emittedLight * colour;
      colour *= material.colour.xyz;
    }
    else {
      incomingLight += Sky(ray) * colour;
      break;
    }
  }
  return incomingLight;
}


void main(){

    ivec2 pixel = ivec2(gl_GlobalInvocationID.xy);
    ivec2 size = imageSize(renderImage);

    if (pixel.x >= size.x || pixel.y >= size.y)
        return;

    vec2 uv = ((vec2(pixel) + 0.5) / vec2(size)) * 2.0 - 1.0;
    int rngState = int(pixel.x * int(1973) + pixel.y * int(9277) + iFrame * int(26699)) | int(1);

    vec4 rayClip = vec4(uv, -1.0, 1.0);

    // clip → view
    vec4 rayView = invProjMat * rayClip;
    rayView /= rayView.w;

    // view → world (direction)
    vec3 rayDir = normalize((invViewMat* vec4(rayView.xyz, 0.0)).xyz);

    Ray ray;
    ray.origin = camPos;
    ray.dir = rayDir;

    vec3 totalIncomingLight = vec3(0);

    for(int rayIndex = 0; rayIndex < NUM_RAYS_PER_PIXEL; rayIndex++){
      totalIncomingLight += TraceRay(ray, rngState);
    }

    vec3 pixelColour = totalIncomingLight / NUM_RAYS_PER_PIXEL;

    vec3 previousColour = imageLoad(accumulation, pixel).rgb;
    previousColour += pixelColour;

    imageStore(accumulation, pixel, vec4(previousColour, 1.0));

    vec3 accumulatedColour = previousColour / iFrame;
    imageStore(renderImage, pixel, vec4(accumulatedColour, 1.0));

}
