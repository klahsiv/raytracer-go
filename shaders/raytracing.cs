#version 460 core

#define PI 3.141592653589793 
#define MAX_SPHERE_COUNT 12
#define MAX_BOUNCE_COUNT 8 
#define NUM_RAYS_PER_PIXEL 10

layout(local_size_x = 8, local_size_y = 8) in;

layout(rgba32f, binding = 0) uniform image2D renderImage;

uniform mat4 invViewMat;
uniform mat4 invProjMat;
uniform vec3 camPos;

uniform int iFrame;
uniform int sphereCount;
uniform int triangleCount;

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
layout(std430, binding = 1) buffer Spheres {
  Sphere spheres[];
};
layout(std430, binding = 2) buffer Materials{
  RayTracingMaterial materials[];
};
layout(std430, binding = 3) buffer Triangles{
  Triangle triangles[];
};
layout(std430, binding = 4) buffer Accumulation{
  vec4 accumulation[];
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
HitInfo RaySphere(Ray ray, vec3 sphereCenter, float sphereRadius){

  HitInfo hitInfo;
  hitInfo.didHit = 0;

  vec3 oc = ray.origin - sphereCenter;

  float a = dot(ray.dir, ray.dir);
  float b = 2.0 * dot(oc, ray.dir);
  float c = dot(oc, oc) - sphereRadius * sphereRadius;

  float discriminant = b * b - 4.0 * a * c;

  if(discriminant >= 0.0){
    float dst = (-b - sqrt(discriminant)) / (2.0 * a);

    if(dst >= 0.0){
      hitInfo.didHit = 1;
      hitInfo.dst = dst;
      hitInfo.hitPoint = ray.origin + ray.dir * dst;
      hitInfo.normal = normalize(hitInfo.hitPoint - sphereCenter);
    }
  }

  return hitInfo;
}

HitInfo RayTriangle(Ray ray, Triangle tri){

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



  HitInfo hitInfo;
  hitInfo.didHit = (determinant >= 1E-6 && dst >= 0 && u >= 0 && v >= 0 && w >= 0) ? 1 : 0;
  hitInfo.hitPoint = ray.origin + ray.dir * dst;
  hitInfo.normal = normalize(tri.normalA * w + tri.normalB * u + tri.normalC * v).xyz;
  hitInfo.dst = dst;

  return hitInfo;
}

HitInfo CalculateRayCollison(Ray ray){

  HitInfo closestHit;
  closestHit.didHit = 0;
  closestHit.dst = 10000.0;
  RayTracingMaterial initMaterial = {vec4(vec3(0.0), 1.0), vec4(0.0)};
  closestHit.material = initMaterial;

  for(int i = 0; i < int(sphereCount); i++){
    vec3 center = spheres[i].pos.xyz;
    float radius = spheres[i].pos.w;
    int materialIdx = int(spheres[i].material.x);
    RayTracingMaterial material = materials[materialIdx];

    HitInfo hitInfo = RaySphere(ray, center, radius);

    if(hitInfo.didHit > 0 && hitInfo.dst < closestHit.dst) {
      closestHit = hitInfo;
      closestHit.material = material;
    }
  }

  /*
  for(int i = 0; i < int(triangleCount); i++){
    Triangle tri = triangles[i];
    int materialIdx = int(spheres[i].material.x);
    RayTracingMaterial material = materials[materialIdx];

    HitInfo hitInfo = RayTriangle(ray, tri);

    if(hitInfo.didHit > 0 && hitInfo.dst < closestHit.dst) {
      closestHit = hitInfo;
      closestHit.material = material;
    }
  }
  */
  return closestHit;
}

vec3 TraceRay(Ray ray, inout int rngState){
  vec3 colour = vec3(1.0);
  vec3 incomingLight = vec3(0.0);
  for(int i = 0; i < MAX_BOUNCE_COUNT; i++){
    HitInfo hitInfo = CalculateRayCollison(ray);

    if(hitInfo.didHit > 0){
      ray.origin = hitInfo.hitPoint;
      ray.dir = normalize(hitInfo.normal + RandomUnitVector(rngState));

      RayTracingMaterial material = hitInfo.material;
      vec3 emittedLight = material.emission.xyz * material.emission.w;
      incomingLight += emittedLight * colour;
      colour *= material.colour.xyz;
    }
    else {
      //incomingLight += vec3(0.2, 0.3, 0.5);
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

    /*
    int idx = int(gl_FragCoord.x) + int(gl_FragCoord.y * resolution.x);
    accumulation[idx].xyz += pixelColor;
    vec3 texel = accumulation[idx].xyz / iFrame;
    */

    imageStore(renderImage, pixel, vec4(pixelColour, 1.0));

}
