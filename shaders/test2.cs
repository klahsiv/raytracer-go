#version 460 core

#define PI 3.141592653589793
#define MAX_BOUNCE_COUNT 8
#define NUM_RAYS_PER_PIXEL 100

layout(local_size_x = 8, local_size_y = 8) in;

layout(rgba32f, binding = 0) uniform image2D renderImage;

uniform vec2 resolution;

uniform mat4 invViewMat;
uniform mat4 invProjMat;
uniform vec3 camPos;

uniform int iFrame;
uniform int sphereCount;
uniform int triangleCount;

struct Sphere {
    vec4 pos;
    vec4 material;
};

struct Triangle {
    vec4 posA, posB, posC;
    vec4 normalA, normalB, normalC;
    vec4 material;
};

struct RayTracingMaterial {
    vec4 colour;
    vec4 emission;
};

layout(std430, binding = 1) buffer Spheres {
    Sphere spheres[];
};

layout(std430, binding = 2) buffer Materials {
    RayTracingMaterial materials[];
};

layout(std430, binding = 3) buffer Triangles {
    Triangle triangles[];
};

layout(std430, binding = 4) buffer Accumulation {
    vec4 accumulation[];
};

struct Ray {
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
    float z = RandomFloat01(state) * 2.0 - 1.0;
    float a = RandomFloat01(state) * 2.0 * PI;
    float r = sqrt(1.0 - z * z);

    return vec3(
        r * cos(a),
        r * sin(a),
        z
    );
}

HitInfo RaySphere(Ray ray, vec3 center, float radius)
{
    HitInfo hit;
    hit.didHit = 0;

    vec3 oc = ray.origin - center;

    float a = dot(ray.dir, ray.dir);
    float b = 2.0 * dot(oc, ray.dir);
    float c = dot(oc, oc) - radius * radius;

    float d = b*b - 4.0*a*c;

    if(d >= 0.0)
    {
        float t = (-b - sqrt(d)) / (2.0 * a);

        if(t >= 0.0)
        {
            hit.didHit = 1;
            hit.dst = t;
            hit.hitPoint = ray.origin + ray.dir * t;
            hit.normal = normalize(hit.hitPoint - center);
        }
    }

    return hit;
}

HitInfo CalculateRayCollision(Ray ray)
{
    HitInfo closest;
    closest.didHit = 0;
    closest.dst = 1e20;

    for(int i = 0; i < sphereCount; i++)
    {
        Sphere sphere = spheres[i];

        HitInfo hit = RaySphere(
            ray,
            sphere.pos.xyz,
            sphere.pos.w
        );

        if(hit.didHit == 1 && hit.dst < closest.dst)
        {
            closest = hit;

            int materialIdx = int(sphere.material.x);
            closest.material = materials[materialIdx];
        }
    }

    return closest;
}

vec3 TraceRay(Ray ray, inout int rngState)
{
    vec3 throughput = vec3(1.0);
    vec3 incomingLight = vec3(0.0);

    for(int bounce = 0; bounce < MAX_BOUNCE_COUNT; bounce++)
    {
        HitInfo hit = CalculateRayCollision(ray);

        if(hit.didHit == 0)
            break;

        RayTracingMaterial material = hit.material;

        vec3 emitted =
            material.emission.xyz *
            material.emission.w;

        incomingLight += emitted * throughput;

        throughput *= material.colour.xyz;

        ray.origin = hit.hitPoint + hit.normal * 0.001;
        ray.dir = normalize(
            hit.normal + RandomUnitVector(rngState)
        );
    }

    return incomingLight;
}


void main(){
    ivec2 pixel = ivec2(gl_GlobalInvocationID.xy);
    ivec2 size = imageSize(renderImage);

    if(pixel.x >= int(size.x) ||
       pixel.y >= int(size.y))
        return;
vec2 uv =
    ((vec2(pixel) + 0.5) / vec2(size)) * 2.0 - 1.0;

    int rngState = int(pixel.x * int(1973) + pixel.y * int(9277) + iFrame * int(26699)) | int(1);

uv.x *= float(size.x) / float(size.y);

// vertical FOV
float fov = radians(45.0);

float scale = tan(fov * 0.5);

vec3 rayDir = normalize(vec3(
    uv.x * scale,
    uv.y * scale,
    -1.0
));

rayDir =
    normalize(
        (invViewMat * vec4(rayDir, 0.0)).xyz
    );

    Ray ray;
    ray.origin = camPos;
    ray.dir = rayDir;

    vec3 totalIncomingLight = vec3(0);

    for(int rayIndex = 0; rayIndex < NUM_RAYS_PER_PIXEL; rayIndex++){
      totalIncomingLight += TraceRay(ray, rngState);
    }

    vec3 pixelColour = totalIncomingLight / NUM_RAYS_PER_PIXEL;
imageStore(
    renderImage,
    pixel,
    vec4(pixelColour, 1.0)
);
}
