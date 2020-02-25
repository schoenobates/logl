package main

import (
    "errors"
    "fmt"
    "github.com/go-gl/gl/v3.3-core/gl"
    "logl/render"
    "os"
    "runtime"
    "strings"
)

// ensure main is on the main thread - see https://github.com/golang/go/wiki/LockOSThread
func init() {
    runtime.LockOSThread()
}

// --------------------------------------------------------------------------------------------------------
// ShaderType
// --------------------------------------------------------------------------------------------------------

// Upgraded shader type constant string with support for printing the type
type ShaderType uint32

func (s ShaderType) String() string {
    switch s {
    case gl.FRAGMENT_SHADER:
        return "Fragment Shader"
    case gl.VERTEX_SHADER:
        return "Vertex Shader"
    default:
        return "Unknown"
    }
}

// Supported shader types
const (
    FragmentShader ShaderType = gl.FRAGMENT_SHADER
    VertexShader   ShaderType = gl.VERTEX_SHADER
)

// --------------------------------------------------------------------------------------------------------
// Shader
// --------------------------------------------------------------------------------------------------------
type Shader struct {
    Type   ShaderType
    Source string
    Ptr    uint32
}

func (s *Shader) Delete() {
    gl.DeleteShader(s.Ptr)
}

func NewShader(shaderType ShaderType, source string) (*Shader, error) {
    shader := gl.CreateShader(uint32(shaderType))

    if shader == 0 {
        // error occurred
        return nil, fmt.Errorf("failed to create shader: type = %s", shaderType)
    }

    csources, free := gl.Strs(source)

    // ensure that we clean up the C pointer
    defer free()

    gl.ShaderSource(shader, 1, csources, nil)
    gl.CompileShader(shader)

    var status int32
    gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)

    if status != gl.FALSE {
        return &Shader{
            Type:   shaderType,
            Source: source,
            Ptr:    shader,
        }, nil
    }

    var logLength int32
    gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

    log := strings.Repeat("\x00", int(logLength+1))
    gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

    return nil, fmt.Errorf(
        "failed to compile shader: type = %s, source = %s, log = %s",
        shaderType, source, log,
    )

}

// --------------------------------------------------------------------------------------------------------
// Program
// --------------------------------------------------------------------------------------------------------
type Program struct {
    Ptr uint32
}

func NewProgram(shaders ...*Shader) (*Program, error) {

    if len(shaders) == 0 {
        return nil, errors.New("no shaders specified to link into program")
    }
    prog := gl.CreateProgram()

    if prog == 0 {
        return nil, fmt.Errorf("failed to create program")
    }

    for _, shader := range shaders {
        gl.AttachShader(prog, shader.Ptr)
    }

    gl.LinkProgram(prog)

    var status int32
    gl.GetProgramiv(prog, gl.LINK_STATUS, &status)

    if status != gl.FALSE {
        return &Program{prog}, nil
    }

    var logLength int32
    gl.GetProgramiv(prog, gl.INFO_LOG_LENGTH, &logLength)

    olog := strings.Repeat("\x00", int(logLength+1))
    gl.GetProgramInfoLog(prog, logLength, nil, gl.Str(olog))

    return nil, fmt.Errorf(
        "failed to link program: log = %s",
        olog,
    )
}

func main() {

    fmt.Println("application starting")

    window, err := render.NewWindow("Simple Triangle", true, 800, 600)
    if err != nil {
        fmt.Printf("failed to create window (and associated gl/glfw resources): error = %s", err)
        os.Exit(1)
    }

    vertices := []float32{
        -0.5, -0.5, 0.0,
        0.5, -0.5, 0.0,
        0.0, 0.5, 0.0,
    }

    vshs := `
#version 330 core
layout (location = 0) in vec3 aPos;

void main()
{
    gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
} 
` + "\x00"

    fshs := `
#version 330 core
out vec4 FragColor;

void main()
{
    FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
}
` + "\x00"

    // setup the program work
    vsh, err := NewShader(VertexShader, vshs)
    if err != nil {
        fmt.Printf("failed to create vertex shader: error = %s", err)
        os.Exit(1)
    }

    fsh, err := NewShader(FragmentShader, fshs)
    if err != nil {
        fmt.Printf("failed to create fragment shader: error = %s", err)
        os.Exit(1)
    }

    program, err := NewProgram(vsh, fsh)

    if err != nil {
        fmt.Printf("failed to link program: error = %s", err)
        os.Exit(1)
    }

    vsh.Delete()
    fsh.Delete()

    // set up the various VBOs and VAOs here before the render
    // buffer to store the data in
    var vbo uint32
    gl.GenBuffers(1, &vbo)

    // vao setup
    var vao uint32
    gl.GenVertexArrays(1, &vao)

    gl.BindVertexArray(vao)

    // bind the buffer
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 12, nil)
    gl.EnableVertexAttribArray(0)

    gl.ClearColor(1.0, 1.0, 1.0, 1.0)
    window.Render(func() {
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

        gl.UseProgram(program.Ptr)
        gl.BindVertexArray(vao)

        // draw vao
        gl.DrawArrays(gl.TRIANGLES, 0, 3)
    })

    window.Destroy()

    fmt.Println("application exiting")
}
