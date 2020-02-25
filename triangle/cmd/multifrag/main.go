package main

import (
    "fmt"
    "github.com/go-gl/gl/v3.3-core/gl"
    "logl/render"
    "os"
)

func main() {
    window, err := render.NewWindow("Multi Fragment", false, 800, 600)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    fshg, err := render.ReadShader(render.FragmentShader, "fragg.glsl")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fshy, err := render.ReadShader(render.FragmentShader, "fragy.glsl")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    vsh, err := render.ReadShader(render.VertexShader, "../vert.glsl")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    progg, err := render.NewProgram(vsh, fshg)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    progy, err := render.NewProgram(vsh, fshy)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    fshg.Delete()
    fshy.Delete()
    vsh.Delete()

    var vao0, vao1 uint32
    var vbo0, vbo1 uint32

    gl.GenVertexArrays(1, &vao0)
    gl.GenVertexArrays(1, &vao1)

    gl.GenBuffers(1, &vbo0)
    gl.GenBuffers(1, &vbo1)

    vert0 := []float32{
        -1.0, 0.0, 0.0,
        -0.5, 0.5, 0.0,
        0.0, 0.0, 0.0,
    }

    vert1 := []float32{
        0.0, 0.0, 0.0,
        0.5, 0.5, 0.0,
        1.0, 0.0, 0.0,
    }

    gl.BindVertexArray(vao0)
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo0)
    gl.BufferData(gl.ARRAY_BUFFER, len(vert0)*4, gl.Ptr(vert0), gl.STATIC_DRAW)
    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 12, nil)
    gl.EnableVertexAttribArray(0)

    gl.BindVertexArray(vao1)
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo1)
    gl.BufferData(gl.ARRAY_BUFFER, len(vert1)*4, gl.Ptr(vert1), gl.STATIC_DRAW)
    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 12, nil)
    gl.EnableVertexAttribArray(0)

    gl.ClearColor(0.0, 0.0, 0.0, 1.0)
    window.Render(func() {
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
        progg.Use()
        gl.BindVertexArray(vao0)
        gl.DrawArrays(gl.TRIANGLES, 0, 3)

        progy.Use()
        gl.BindVertexArray(vao1)
        gl.DrawArrays(gl.TRIANGLES, 0, 3)
    })

}
