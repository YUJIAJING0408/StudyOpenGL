#version 330
in vec2 fragTexCoord;
out vec4 outputColor;
uniform sampler2D tex;
uniform vec4 textColor;
void main()
{
    vec4 sampled = vec4(1.0, 1.0, 1.0, texture(tex, fragTexCoord).r);
    outputColor = textColor * sampled;
}