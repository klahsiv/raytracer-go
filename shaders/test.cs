#version 460 core

layout(local_size_x = 8, local_size_y = 8) in;

layout(rgba32f, binding = 0) uniform image2D renderImage;

uniform mat4 viewMat;
uniform mat4 projMat;
uniform vec3 camPos;

uniform int iFrame;
uniform int sphereCount;
uniform int triangleCount;

struct Sphere{
    vec4 pos;
    vec4 material;
};

struct Triangle{
    vec4 posA, posB, posC;
    vec4 normalA, normalB, normalC;
    vec4 material;
};

struct RayTracingMaterial{
    vec4 colour;
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

bool RaySphere(
    Ray ray,
    vec3 center,
    float radius,
    out float t
){
    vec3 oc = ray.origin - center;

    float a = dot(ray.dir, ray.dir);
    float b = 2.0 * dot(oc, ray.dir);
    float c = dot(oc, oc) - radius * radius;

    float discriminant = b * b - 4.0 * a * c;

    if(discriminant < 0.0)
        return false;

    t = (-b - sqrt(discriminant)) / (2.0 * a);

    return t > 0.0;
}

void main()
{
    ivec2 pixel = ivec2(gl_GlobalInvocationID.xy);
    ivec2 size = imageSize(renderImage);

    if(pixel.x >= size.x || pixel.y >= size.y)
        return;

    // Screen UV
    vec2 uv = ((vec2(pixel) + 0.5) / vec2(size)) * 2.0 - 1.0;

    uv.y *= -1.0;

    // Aspect ratio correction
    uv.x *= float(size.x) / float(size.y);

    // Simple camera ray
    vec3 rayDir = normalize(vec3(uv, -1.0));

    // Apply inverse view matrix
    rayDir = normalize((viewMat * vec4(rayDir, 0.0)).xyz);

    Ray ray;
    ray.origin = camPos;
    //ray.origin = vec3(0.0, 0.0, 15.0);
    ray.dir = rayDir;

    // Hardcoded sphere
    vec3 sphereCenter = spheres[1].pos.xyz;
    float sphereRadius = spheres[1].pos.w;

    float t;

    if(RaySphere(ray, sphereCenter, sphereRadius, t))
    {
        imageStore(renderImage, pixel, vec4(1.0, 0.0, 0.0, 1.0));
    }
    else
    {
        // Sky gradient
        vec3 sky = mix(
            vec3(1.0),
            vec3(0.5, 0.7, 1.0),
            0.5 * (rayDir.y + 1.0)
        );

        imageStore(renderImage, pixel, vec4(sky, 1.0));
    }
}
