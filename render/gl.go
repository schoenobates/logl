package render

import (
    "errors"
    "fmt"
    "github.com/go-gl/gl/v3.3-core/gl"
    "io/ioutil"
    "strings"
)

// basic Vector4 type
type Vec4 [4]float32

const (
    // Vec4 component synonyms - RGBA
    R = iota
    G
    B
    A

    // Vec4 component synonyms - XYZW
    X = iota
    Y
    Z
    W
)

// holds a colour type
type Color struct{ r, b, g, a float32 }

var (
    White = Color{1, 1, 1, 1}
    Black = Color{0, 0, 0, 1}
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
    ptr    uint32
}

func (s *Shader) Delete() {
    gl.DeleteShader(s.ptr)
}

// Will read and add the null terminator to the given shader at the specified path
func ReadShader(shaderType ShaderType, path string) (*Shader, error) {

    contents, err := ioutil.ReadFile(path)

    if err != nil {
        return nil, err
    }

    contents = append(contents, []byte("\x00")...)

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
            ptr:    shader,
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
    ptr uint32
}

// sets this as the active program
func (p *Program) Use() {
    gl.UseProgram(p.ptr)
}

// sets a uniform boolean value specified by the given name - this will set the value as an integer - 1 = true, 0 = false
func (p *Program) Bool(name string, value bool) error {
    if value {
        return p.Integer(name, 1)
    } else {
        return p.Integer(name, 0)
    }
}

// sets an integer uniform value
func (p *Program) Integer(name string, value int32) error {
    if location, err := p.uniform(name); err == nil {
        gl.Uniform1i(location, value)
        return nil
    } else {
        return err
    }
}

// sets a float32 uniform value
func (p *Program) Float(name string, value float32) error {
    if location, err := p.uniform(name); err == nil {
        gl.Uniform1f(location, value)
        return nil
    } else {
        return err
    }
}

// sets a vec4 uniform value
func (p *Program) Vec4(name string, value Vec4) error {
    if location, err:= p.uniform(name); err == nil {
        gl.Uniform4f(location, value[0], value[1], value[2], value[3])
        return nil
    } else {
        return err
    }
}

// gets the uniform location for the given name
func (p *Program) uniform(name string) (int32, error) {
    location := gl.GetUniformLocation(p.ptr, gl.Str(name+"\x00"))

    if location == -1 {
        return 0, errors.New(fmt.Sprintf("failed to locate uniform: name = %s", name))
    }

    return location, nil
}

// Creates a new program instance from the given shader set. Callers to this function are required to manage the
// shader cleanup (Delete) - this is not done here.
func NewProgram(shaders ...*Shader) (*Program, error) {

    if len(shaders) == 0 {
        return nil, errors.New("no shaders specified to link into program")
    }

    prog := gl.CreateProgram()

    if prog == 0 {
        return nil, fmt.Errorf("failed to create program")
    }

    for _, shader := range shaders {
        gl.AttachShader(prog, shader.ptr)
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
