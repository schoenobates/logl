package render

import (
    "errors"
    "fmt"
    "github.com/go-gl/gl/v3.3-core/gl"
    "github.com/go-gl/mathgl/mgl32"
    "image"
    "image/draw"
    _ "image/gif"
    _ "image/jpeg"
    _ "image/png"
    "io/ioutil"
    "os"
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
// Texture
// --------------------------------------------------------------------------------------------------------

// texture wrap options
type TextureWrap uint32

const (
    Repeat         TextureWrap = gl.REPEAT
    ClampToEdge    TextureWrap = gl.CLAMP_TO_EDGE
    MirroredRepeat TextureWrap = gl.MIRRORED_REPEAT
)

// texture filter options for min/mag functions
type TextureFilter uint32

const (
    Nearest              TextureFilter = gl.NEAREST
    Linear               TextureFilter = gl.LINEAR
    NearestMipmapNearest TextureFilter = gl.NEAREST_MIPMAP_NEAREST
    LinearMipmapNearest  TextureFilter = gl.LINEAR_MIPMAP_NEAREST
    NearestMipmapLinear  TextureFilter = gl.NEAREST_MIPMAP_LINEAR
    LinearMipmapLinear   TextureFilter = gl.LINEAR_MIPMAP_LINEAR
)

type TextureUnit struct {
    bind  int
    index int
}

func (t TextureUnit) Index() int32 {
    return int32(t.index)
}

var (
    TextureUnit0  = TextureUnit{gl.TEXTURE0, 0}
    TextureUnit1  = TextureUnit{gl.TEXTURE1, 1}
    TextureUnit2  = TextureUnit{gl.TEXTURE2, 2}
    TextureUnit3  = TextureUnit{gl.TEXTURE3, 3}
    TextureUnit4  = TextureUnit{gl.TEXTURE4, 4}
    TextureUnit5  = TextureUnit{gl.TEXTURE5, 5}
    TextureUnit6  = TextureUnit{gl.TEXTURE6, 6}
    TextureUnit7  = TextureUnit{gl.TEXTURE7, 7}
    TextureUnit8  = TextureUnit{gl.TEXTURE8, 8}
    TextureUnit9  = TextureUnit{gl.TEXTURE9, 9}
    TextureUnit10 = TextureUnit{gl.TEXTURE10, 10}
    TextureUnit11 = TextureUnit{gl.TEXTURE11, 11}
    TextureUnit12 = TextureUnit{gl.TEXTURE12, 12}
    TextureUnit13 = TextureUnit{gl.TEXTURE13, 13}
    TextureUnit14 = TextureUnit{gl.TEXTURE14, 14}
    TextureUnit15 = TextureUnit{gl.TEXTURE15, 15}
)

// holds the actual texture reference
type Texture struct {
    ptr uint32
}

// binds this texture for usage to the given texture uint
func (t *Texture) Bind(textureUnit TextureUnit) {
    gl.ActiveTexture(uint32(textureUnit.bind))
    gl.BindTexture(gl.TEXTURE_2D, t.ptr)
}

// options used when creating the texture
type TextureOpts struct {
    GenMipMap bool
    WrapS     TextureWrap
    WrapT     TextureWrap
    MinFilter TextureFilter
    MagFilter TextureFilter
    FlipY     bool
}

// Reads a texture from the given path
func ReadTexture(path string, opts TextureOpts) (*Texture, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }

    img, _, err := image.Decode(f)

    if err != nil {
        return nil, err
    }

    rgba := image.NewRGBA(img.Bounds())
    if rgba.Stride != rgba.Rect.Size().X*4 {
        return nil, fmt.Errorf("unsupported stride")
    }
    draw.Draw(rgba, rgba.Bounds(), img, image.Point{X: 0, Y: 0}, draw.Src)

    if opts.FlipY {
        for y := 0; y < rgba.Bounds().Size().Y / 2; y++ {
            for x := 0; x < rgba.Bounds().Size().X; x++ {
                c0 := rgba.At(x, y)
                c1 := rgba.At(x, rgba.Bounds().Size().Y - 1 - y)

                rgba.Set(x, y, c1)
                rgba.Set(x, rgba.Bounds().Size().Y - 1 - y, c0)
            }
        }
    }

    var texture uint32
    gl.GenTextures(1, &texture)
    gl.BindTexture(gl.TEXTURE_2D, texture)

    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, int32(opts.WrapS))
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, int32(opts.WrapT))
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, int32(opts.MinFilter))
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, int32(opts.MagFilter))

    gl.TexImage2D(
        gl.TEXTURE_2D,
        0,
        gl.RGBA,
        int32(rgba.Rect.Size().X),
        int32(rgba.Rect.Size().Y),
        0,
        gl.RGBA,
        gl.UNSIGNED_BYTE,
        gl.Ptr(rgba.Pix),
    )

    if opts.GenMipMap {
        gl.GenerateMipmap(gl.TEXTURE_2D)
    }

    return &Texture{texture}, nil

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
    if location, err := p.uniform(name); err == nil {
        gl.Uniform4f(location, value[0], value[1], value[2], value[3])
        return nil
    } else {
        return err
    }
}

func (p *Program) Mat4(name string, value mgl32.Mat4) error {
    if location, err := p.uniform(name); err == nil {
        gl.UniformMatrix4fv(location, 1, false, &value[0])
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
