#version 330 core
out vec4 FragColor;

in vec3 ourColor;
in vec2 texCoord;

uniform bool ismix;
uniform bool isblend;
uniform float mixture;

uniform sampler2D containerTexture;
uniform sampler2D awesomeTexture;

void main()
{
    if (ismix) {
        FragColor = texture(containerTexture, texCoord) * vec4(ourColor, 1.0);
    } else if (isblend) {
        FragColor = mix(texture(containerTexture, texCoord), texture(awesomeTexture, texCoord), mixture);
    } else {
        FragColor = texture(containerTexture, texCoord);
    }
}