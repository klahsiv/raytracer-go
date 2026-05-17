#version 460 core

layout(local_size_x = 8, local_size_y = 8) in;

layout(rgba8, binding = 0) uniform image2D renderImage;

void main()
{
    ivec2 pixel = ivec2(gl_GlobalInvocationID.xy);
    ivec2 size = imageSize(renderImage);

    if (pixel.x >= size.x || pixel.y >= size.y)
        return;

    vec2 uv = vec2(pixel) / vec2(size);

    vec3 color = vec3(uv, 0.0);

    imageStore(renderImage, pixel, vec4(color, 1.0));
}
