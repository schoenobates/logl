package main

import (
    "fmt"
    "github.com/go-gl/gl/v3.3-core/gl"
    "github.com/go-gl/glfw/v3.3/glfw"
    "logl/render"
    "math"
    "os"
)

func main() {
    window, err := render.NewWindow("Uniform", false, 800, 600)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    defer window.Destroy()

    var fsh, vsh *render.Shader
    if fsh, err = render.ReadShader(render.FragmentShader, "frag.glsl"); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    if vsh, err = render.ReadShader(render.VertexShader, "vert.glsl"); err != nil {
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

    vertices := []float32{
        -0.5, -0.5, 0,
        0, 0.5, 0,
        0.5, -0.5, 0,
    }

    var vao uint32
    gl.GenVertexArrays(1, &vao)

    var vbo uint32
    gl.GenBuffers(1, &vbo)

    // setup the vao object
    gl.BindVertexArray(vao)

    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 12, nil)
    gl.EnableVertexAttribArray(0)

    window.ClearColor(render.Black)

    c := render.Vec4{0.0, 0.0, 0.0, 1}

    window.Render(func() {
        if window.IsPressed(glfw.KeyEscape) {
            window.Close()
        }

        prog.Use()

        c[render.G] = float32((math.Sin(window.Time()) / 2.0) + 0.5)

        err = prog.Vec4("ourColor", c)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        gl.BindVertexArray(vao)
        gl.DrawArrays(gl.TRIANGLES, 0, 3)
    })

}
