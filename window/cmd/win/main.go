package main

// --------------------------------------------------------------------------------------------------------------------
// Source for the tutorial https://learnopengl.com/Getting-started/Hello-Window
// --------------------------------------------------------------------------------------------------------------------

import (
    "fmt"
    "github.com/go-gl/gl/v3.3-core/gl"
    "github.com/go-gl/glfw/v3.3/glfw"
    "math"
    "os"
    "runtime"
)

func init() {
    runtime.LockOSThread()
}

func main() {

    // ----------------------------------------------------------------------------------------------------
    // main initialisation blocks
    // ----------------------------------------------------------------------------------------------------
    if err := gl.Init(); err != nil {
        fmt.Printf("failed to initialise GLFW: error = %s\n", err)
        os.Exit(1)
    }

    if err := glfw.Init(); err != nil {
        fmt.Printf("failed to initialise GLFW: error = %s", err)
        os.Exit(1)
    }

    defer glfw.Terminate()

    if err := glfw.Init(); err != nil {
        fmt.Printf("failed to initialise GLFW: error = %s", err)
        os.Exit(1)
    }

    defer glfw.Terminate()

    glfw.WindowHint(glfw.ContextVersionMajor, 3)
    glfw.WindowHint(glfw.ContextVersionMinor, 3)
    glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

    // MacOS compat
    glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

    window, err := glfw.CreateWindow(800, 600, "LearnOpenGL", nil, nil)

    if err != nil {
        fmt.Printf("failed to create new window instance: error = %s\n", err)
        os.Exit(1)
    }

    window.MakeContextCurrent()

    version := gl.GoStr(gl.GetString(gl.VERSION))
    fmt.Printf("running with opengl: version = %s\n", version)


    // handle window resize functions
    window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
        gl.Viewport(0, 0, int32(width), int32(height))
    })

    cycle := false
    colour := float32(0.0)
    cv := math.Sin(float64(colour))
    // main render loop
    for !window.ShouldClose() {
        // handle input
        if window.GetKey(glfw.KeyEscape) == glfw.Press && !window.ShouldClose() {
            fmt.Println("shutting window to close on ESC key press")
            window.SetShouldClose(true)
        }

        if window.GetKey(glfw.KeyA) == glfw.Press {
            cycle = true
        }

        if window.GetKey(glfw.KeyS) == glfw.Press {
            cycle = false
        }

        // render
        if cycle {
            colour += 0.0001
            cv = math.Sin(float64(colour))
            gl.ClearColor(float32(cv), float32(cv), float32(cv), 1.0)
        } else {
            gl.ClearColor(0.2, 0.3, 0.3, 1.0)
        }
        gl.Clear(gl.COLOR_BUFFER_BIT)

        // swap and poll
        window.SwapBuffers()
        glfw.PollEvents()
    }
}

