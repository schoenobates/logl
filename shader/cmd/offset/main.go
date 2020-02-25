package main

import (
    "fmt"
    "github.com/go-gl/gl/v3.3-core/gl"
    "logl/render"
    "os"
)

func main() {
    window, err := render.NewWindow("More Attributes!", false, 800, 600)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    vertices := []float32 {
        // positions      // colors
         0.5, -0.5, 0.0,  1.0, 0.0, 0.0,   // bottom right
        -0.5, -0.5, 0.0,  0.0, 1.0, 0.0,   // bottom left
         0.0,  0.5, 0.0,  0.0, 0.0, 1.0,    // top
    }

    fsh, err := render.ReadShader(render.FragmentShader, "frag.glsl")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    vsh, err := render.ReadShader(render.VertexShader, "vert.glsl")
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

    window.ClearColor(render.White)

    var vao, vbo uint32
    gl.GenVertexArrays(1, &vao)
    gl.GenBuffers(1, &vbo)

    gl.BindVertexArray(vao)
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.BufferData(gl.ARRAY_BUFFER, len(vertices) * 4, gl.Ptr(vertices), gl.STATIC_DRAW)

    // position attribute
    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6 * 4, nil)
    gl.EnableVertexAttribArray(0)

    // colour attribute
    gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6 * 4, gl.PtrOffset(3*4))
    gl.EnableVertexAttribArray(1)

    window.Render(func() {
        prog.Use()
        err = prog.Float("xoffset", -0.5)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        gl.BindVertexArray(vao)
        gl.DrawArrays(gl.TRIANGLES, 0, 6)
    })


}
