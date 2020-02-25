package main

import (
    "fmt"
    "github.com/go-gl/gl/v3.3-core/gl"
    "logl/render"
    "os"
    "runtime"
)

func init() {
    runtime.LockOSThread()
}

func main() {
    window, err := render.NewWindow("Two Triangles, Multi VBO/VAO", false, 800, 600)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    vsh, err := render.ReadShader(render.VertexShader, "../vert.glsl")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fsh, err := render.ReadShader(render.FragmentShader, "../frag.glsl")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    prog, err := render.NewProgram(vsh, fsh)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    vsh.Delete()
    fsh.Delete()

    var vbo0, vbo1 uint32
    gl.GenBuffers(1, &vbo0)
    gl.GenBuffers(1, &vbo1)

    var vao0, vao1 uint32
    gl.GenVertexArrays(1, &vao0)
    gl.GenVertexArrays(1, &vao1)

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

    gl.ClearColor(1.0, 1.0, 1.0, 1.0)
    window.Render(func() {
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
        gl.UseProgram(prog.Ptr)
        gl.BindVertexArray(vao0)
        gl.DrawArrays(gl.TRIANGLES, 0, 3)
        gl.BindVertexArray(vao1)
        gl.DrawArrays(gl.TRIANGLES, 0, 3)
    })
}
