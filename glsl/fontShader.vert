#version 330
//vertex position
in vec2 vert;
//pass through to fragTexCoord
in vec2 vertTexCoord;
//window res
uniform vec2 resolution;
//pass to frag
out vec2 fragTexCoord;

void main() {
    // 归一化
    vec2 zeroToOne = vert / resolution;
    // [0,1]-> [-1,1]
    vec2 clipSpace = zeroToOne * 2.0 - 1.0;
    fragTexCoord = vertTexCoord;
    vec2 outer = clipSpace * vec2(1, -1);
    gl_Position = vec4(outer, 0, 1);
}