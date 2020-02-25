package render

import (
    "errors"
    "fmt"
    "github.com/go-gl/gl/v3.3-core/gl"
    "io/ioutil"
    "strings"
)

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

// Will read and add the null terminator to the given shader at the specified path
func ReadShader(shaderType ShaderType, path string) (*Shader, error) {

    contents, err := ioutil.ReadFile(path)

    if err != nil {
        return nil, err
    }

    contents = append(contents, []byte("\x00") ...)

    return NewShader(shaderType, string(contents))
}

// Creates a new shader from the specified source string
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

// sets this as the active program
func (p *Program) Use() {
    gl.UseProgram(p.Ptr)
}

// sets a uniform boolean value specified by the given name - this will set the value as an integer - 1 = true, 0 = false
func (p *Program) Bool(name string, value bool) error {
    location := gl.GetUniformLocation(p.Ptr, gl.Str(name))

    if location == -1 {
        return errors.New(fmt.Sprintf("failed to locate uniform: name = %s", name))
    }

    if value {
        gl.Uniform1i(location, 1)
    } else {
        gl.Uniform1i(location, 0)
    }

    return nil
}

func NewProgram(shaders ... *Shader) (*Program, error) {

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
