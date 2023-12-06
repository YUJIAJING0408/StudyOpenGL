#version 330 core
in vec3 Normal;
in vec3 FragPos;

out vec4 color;

struct Material{
    vec3 ambient;
    vec3 diffuse;
    vec3 specular;
    float shininess;
};
struct Light{
    vec3 position;
    vec3 ambient;
    vec3 diffuse;
    vec3 specular;
};

uniform Material material;
uniform Light light;
//uniform vec3 lightColor;
//uniform vec3 lightPos;
uniform vec3 viewPos;

void main()
{
    //环境光
    vec3 ambient = light.ambient * material.ambient;
    
    //漫反射
    vec3 norm = normalize(Normal);//片段法向量单位化
    vec3 lightDir = normalize(light.position - FragPos);//光源射向片段的单位化
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = light.diffuse * (diff * material.diffuse);

    //镜面反射
    vec3 viewDir = normalize(viewPos - FragPos);
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0),  material.shininess);
    vec3 specular = light.specular * (spec * material.specular);

    //整合渲染
    vec3 result = (specular +ambient + diffuse);
    color = vec4(result, 1.0f);
}