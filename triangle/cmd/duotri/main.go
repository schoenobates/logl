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
    window, err := render.NewWindow("Two Triangles Single VBO", true, 800, 600)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    defer window.Destroy()

    vertices := []float32{
        -1.0, 0.0, 0.0,
        -0.5, 0.5, 0.0,
        0.0, 0.0, 0.0,
        0.0, 0.0, 0.0,
        0.5, 0.5, 0.0,
        1.0, 0.0, 0.0,
    }

    // make the shaders and program
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

    // var vbo uint32
    // gl.GenBuffers(1, &vbo)
    //
    // // vao setup
    // var vao uint32
    // gl.GenVertexArrays(1, &vao)

    var vbo uint32
    gl.GenBuffers(1, &vbo)

    // set up the buffers
    var vao uint32
    gl.GenVertexArrays(1, &vao)

    // setup the vao object
    gl.BindVertexArray(vao)

    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 12, nil)
    gl.EnableVertexAttribArray(0)

    gl.ClearColor(.0, .0, .0, 1)
    window.Render(func() {
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
        gl.UseProgram(prog.Ptr)
        gl.BindVertexArray(vao)

        // draw vao
        gl.DrawArrays(gl.TRIANGLES, 0, 6)
    })

}
