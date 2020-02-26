package main

import (
    "fmt"
    "github.com/go-gl/gl/v3.3-core/gl"
    "github.com/go-gl/glfw/v3.3/glfw"
    "logl/render"
    "os"
)

func main() {
    window, err := render.NewWindow("Container Texture", false, 800, 600)
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

    ctexture, err := render.ReadTexture("container.jpg", render.TextureOpts{
        GenMipMap: true,
        WrapS: render.Repeat,  // ClampToEdge ex2
        WrapT: render.Repeat,  // ClampToEdge ex2
        MinFilter: render.Linear,
        MagFilter: render.Linear,
        FlipY: false,
    })

    atexture, err := render.ReadTexture("awesomeface.png", render.TextureOpts{
        GenMipMap: true,
        WrapS: render.Repeat,
        WrapT: render.Repeat,
        MinFilter: render.Linear,
        MagFilter: render.Linear,
        FlipY: true,
    })

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    vertices := []float32{
        // positions       // colors       // texture coords
        0.5,  0.5, 0.0,   1.0, 0.0, 0.0,   0.0, 1.0,   // top right
        0.5, -0.5, 0.0,   0.0, 1.0, 0.0,   1.0, 0.0,   // bottom right
        -0.5, -0.5, 0.0,   0.0, 0.0, 1.0,   0.0, 0.0,  // bottom left
        -0.5,  0.5, 0.0,   1.0, 1.0, 0.0,   0.0, 1.0,  // top left
    }

    /*
    // exercise 2
    vertices := []float32{
       // positions       // colors       // texture coords
       0.5,  0.5, 0.0,   1.0, 0.0, 0.0,   2.0, 2.0,   // top right
       0.5, -0.5, 0.0,   0.0, 1.0, 0.0,   2.0, 0.0,   // bottom right
       -0.5, -0.5, 0.0,   0.0, 0.0, 1.0,   0.0, 0.0,  // bottom left
       -0.5,  0.5, 0.0,   1.0, 1.0, 0.0,   0.0, 2.0,  // top left
   }
     */

    /*
    // exercise 3
    vertices := []float32{
       // positions       // colors       // texture coords
       0.5,  0.5, 0.0,   1.0, 0.0, 0.0,   0.75, .75,   // top right     // ex2 2.0, 2.0
       0.5, -0.5, 0.0,   0.0, 1.0, 0.0,   .75, .25,   // bottom right  // ex2 2.0, 0.0
       -0.5, -0.5, 0.0,   0.0, 0.0, 1.0,   .25, .25,  // bottom left   // ex2 0.0, 0.0
       -0.5,  0.5, 0.0,   1.0, 1.0, 0.0,   0.25, 0.75,  // top left      // ex2 0.0, 2.0
   }
     */

    elements := []uint32 {
        0, 1, 3,
        1, 2, 3,
    }

    var vao, vbo, ebo uint32
    gl.GenVertexArrays(1, &vao)

    gl.GenBuffers(1, &vbo)  // vertex
    gl.GenBuffers(1, &ebo)  // element order

    gl.BindVertexArray(vao)

    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.BufferData(gl.ARRAY_BUFFER, len(vertices) * 4, gl.Ptr(vertices), gl.STATIC_DRAW)

    gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
    gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(elements) * 4, gl.Ptr(elements), gl.STATIC_DRAW)

    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 32, gl.PtrOffset(0))
    gl.EnableVertexAttribArray(0)
    gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 32, gl.PtrOffset(3 * 4))
    gl.EnableVertexAttribArray(1)
    gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 32, gl.PtrOffset(6 * 4))
    gl.EnableVertexAttribArray(2)

    window.ClearColor(render.White)

    mix := false
    blend := false

    var mixture float32 = 0.2

    // texture uniforms - bind prog before use - one time only needed
    prog.Use()
    if err = prog.Integer("containerTexture", render.TextureUnit0.Index()); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    if err = prog.Integer("awesomeTexture", render.TextureUnit1.Index()); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    window.Render(func() {

        if window.IsPressed(glfw.KeyM) {
            mix = true
            blend = false
        }

        if window.IsPressed(glfw.KeyN) {
            mix = false
            blend = false
        }

        if window.IsPressed(glfw.KeyB) {
            mix = false
            blend = true
        }

        // use very small increments for mixture as glfw seems to think the key is pressed continuously - need
        // to handle this in a better way
        if window.IsPressed(glfw.KeyUp) {
            mixture += 0.001

            if mixture > 1 {
                mixture = 1.0
            }
        }

        if window.IsPressed(glfw.KeyDown) {
            mixture -= 0.001

            if mixture < 0 {
                mixture = 0
            }
        }

        prog.Use()

        if err := prog.Float("mixture", mixture); err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        if err := prog.Bool("ismix", mix); err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        if err := prog.Bool("isblend", blend); err != nil {
            fmt.Println(err)
            os.Exit(1)
        }


        ctexture.Bind(render.TextureUnit0)
        atexture.Bind(render.TextureUnit1)

        gl.BindVertexArray(vao)
        gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
    })


}
