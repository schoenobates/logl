package main

import (
    "fmt"
    "github.com/go-gl/gl/v3.3-core/gl"
    "github.com/go-gl/mathgl/mgl32"
    "logl/render"
    "os"
)

func main() {

    /**
      glm::vec4 vec(1.0f, 0.0f, 0.0f, 1.0f);
      glm::mat4 trans = glm::mat4(1.0f);
      trans = glm::translate(trans, glm::vec3(1.0f, 1.0f, 0.0f));
      vec = trans * vec;
      std::cout << vec.x << vec.y << vec.z << std::endl;
    */
    vec := mgl32.Vec4{1.0, 0.0, 0.0, 1.0}
    trans := mgl32.Translate3D(1, 1, 0)
    vec = trans.Mul4x1(vec)
    fmt.Printf("[x=%f, y=%f, z=%f]", vec.X(), vec.Y(), vec.Z())

    window, err := render.NewWindow("Transform", false, 800, 600)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    vsh, err := render.ReadShader(render.VertexShader, "vert.glsl")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fsh, err := render.ReadShader(render.FragmentShader, "frag.glsl")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    prog, err := render.NewProgram(vsh, fsh)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    fsh.Delete()
    vsh.Delete()

    atexture, err := render.ReadTexture("awesomeface.png", render.TextureOpts{
        GenMipMap: true,
        WrapS:     render.Repeat,
        WrapT:     render.Repeat,
        MinFilter: render.Linear,
        MagFilter: render.Linear,
        FlipY:     true,
    })

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    vertices := []float32{
        // positions        // colors       // texture coords
        0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, // top right
        0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, // bottom right
        -0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom left
        -0.5, 0.5, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, // top left
    }

    elements := []uint32{
        0, 1, 3,
        1, 2, 3,
    }

    var vao, vbo, ebo uint32
    gl.GenVertexArrays(1, &vao)

    gl.GenBuffers(1, &vbo) // vertex
    gl.GenBuffers(1, &ebo) // element order

    gl.BindVertexArray(vao)

    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

    gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
    gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(elements)*4, gl.Ptr(elements), gl.STATIC_DRAW)

    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 32, gl.PtrOffset(0))
    gl.EnableVertexAttribArray(0)
    gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 32, gl.PtrOffset(3*4))
    gl.EnableVertexAttribArray(1)
    gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 32, gl.PtrOffset(6*4))
    gl.EnableVertexAttribArray(2)

    window.ClearColor(render.White)

    // trans = mgl32.HomogRotate3DZ(mgl32.DegToRad(90))
    // trans = trans.Mul4(mgl32.Scale3D(0.5, 0.5, 0.5))

    // texture uniforms - bind prog before use - one time only needed
    prog.Use()
    if err = prog.Integer("awesomeTexture", render.TextureUnit0.Index()); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    window.Render(func() {

        prog.Use()
        atexture.Bind(render.TextureUnit0)
        gl.BindVertexArray(vao)

        // rotate
        // trans0 := mgl32.Translate3D(0.5, -0.5, 0)
        // trans0 = trans0.Mul4(mgl32.HomogRotate3DZ(float32(window.Time())))

        // flip around for fun and profit
        // trans0 := mgl32.HomogRotate3DZ(float32(window.Time()))
        // trans0 = trans0.Mul4(mgl32.Translate3D(0.5, -0.5, 0))
        if err := prog.Mat4("transform", trans0); err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

        // scale - ex2
        // scale := float32(math.Sin(window.Time()))
        // trans1 := mgl32.Translate3D(-0.5, 0.5, 0)
        // trans1 = trans1.Mul4(mgl32.Scale3D(scale, scale, scale))
        // if err := prog.Mat4("transform", trans1); err != nil {
        //     fmt.Println(err)
        //     os.Exit(1)
        // }
        // gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
    })

}
