package main

import (
    "fmt"
    "github.com/go-gl/gl/v3.3-core/gl"
    "github.com/go-gl/glfw/v3.3/glfw"
    "logl/render"
    "os"
    "runtime"
)

func init() {
    runtime.LockOSThread()
}

func main() {
    window, err := render.NewWindow("EBO", true, 800, 600)

    if err != nil {
        fmt.Printf("failed to create new window: error = %s", err)
        os.Exit(1)
    }

    defer window.Destroy()

    vertices := []float32{
        0.5, 0.5, 0.0,
        -0.5, 0.5, 0.0,
        -0.5, -0.5, 0.0,
        0.5, -0.5, 0.0,
    }

    indices := []uint32{
        0, 1, 3,
        1, 2, 3,
    }

    vsh, err := render.ReadShader(render.VertexShader, "../vert.glsl")
    if err != nil {
        fmt.Printf("failed to create new vsh: error = %s", err)
        os.Exit(1)
    }

    fsh, err := render.ReadShader(render.FragmentShader, "../frag.glsl")
    if err != nil {
        fmt.Printf("failed to create new fsh: error = %s", err)
        os.Exit(1)
    }

    program, err := render.NewProgram(vsh, fsh)

    if err != nil {
        fmt.Printf("failed to create new program: error = %s", err)
        os.Exit(1)
    }

    fsh.Delete()
    vsh.Delete()

    var ebo uint32
    gl.GenBuffers(1, &ebo)

    var vbo uint32
    gl.GenBuffers(1, &vbo)

    var vao uint32
    gl.GenVertexArrays(1, &vao)

    gl.BindVertexArray(vao)

    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

    gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
    gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 12, nil)
    gl.EnableVertexAttribArray(0)

    gl.ClearColor(.0, .0, .0, 1.0)

    wireframe := false

    window.Render(func() {

        if window.IsPressed(glfw.KeyEscape) {
            window.Close()
        }

        if window.IsPressed(glfw.KeyW) {
            wireframe = !wireframe
        }

        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

        if wireframe {
            gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
        } else {
            gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
        }

        gl.UseProgram(program.Ptr)
        gl.BindVertexArray(vao)
        gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

    })
}
