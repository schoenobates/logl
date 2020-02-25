package render

import (
    "fmt"
    "github.com/go-gl/gl/v3.3-core/gl"
    "github.com/go-gl/glfw/v3.3/glfw"
    "runtime"
)

// the renderer function used to handle the actual drawing
type Renderer func()

// holds the gl stuff together and will call the renderer repeatedly
type Window struct {
    Width int32
    Height int32
    win *glfw.Window
}

// closes the window
func (w *Window) Close() {
    if w.win.ShouldClose() {
        return
    }

    w.win.SetShouldClose(true)
}

// destroys the window and terminates the glfw instance
func (w *Window) Destroy() {
    w.win.Destroy()
    glfw.Terminate()
}

func (w *Window) IsReleased(key glfw.Key) bool {
    return w.win.GetKey(key) == glfw.Release
}

func (w *Window) IsPressed(key glfw.Key) bool {
    return w.win.GetKey(key) == glfw.Press
}

// enters the render loop and will block the caller until exit
func (w *Window) Render(render Renderer) {
    for !w.win.ShouldClose() {

        // input - keyboard, mouse etc
        // defaults for now


        // render
        render()

        // update
        // swap and poll
        w.win.SwapBuffers()
        glfw.PollEvents()
    }
}

// creates a new application window and initialises the GLFW and GL subsystems. This call will also lock the current
// go routine to the OS thread (using runtime.LockOSThread())
func NewWindow(title string, resizable bool, width, height int32) (*Window, error) {

    // ----------------------------------------------------------------------------------------------------
    // main initialisation of glfw
    // ----------------------------------------------------------------------------------------------------
    if err := glfw.Init(); err != nil {
        fmt.Printf("failed to initialise GLFW: error = %s", err)
        return nil, err
    }

    glfw.WindowHint(glfw.ContextVersionMajor, 3)
    glfw.WindowHint(glfw.ContextVersionMinor, 3)
    glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

    if resizable {
        glfw.WindowHint(glfw.Resizable, gl.TRUE)
    } else {
        glfw.WindowHint(glfw.Resizable, gl.FALSE)
    }

    if runtime.GOOS == "darwin" {
        glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
    }

    if width < 1 {
        width = 800
    }

    if height < 1 {
        height = 600
    }

    window, err := glfw.CreateWindow(int(width), int(height), title, nil, nil)

    if err != nil {
        return nil, err
    }

    window.MakeContextCurrent()

    if err := gl.Init(); err != nil {
        fmt.Printf("failed to initialise GLFW: error = %s\n", err)
        return nil, err
    }

    version := gl.GoStr(gl.GetString(gl.VERSION))
    fmt.Printf("running with opengl: version = %s\n", version)

    win := &Window{
        win: window,
        Width: width,
        Height: height,
    }

    // custom handling of the window size changes
    window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
        gl.Viewport(0, 0, int32(width), int32(height))
        win.Width = int32(width)
        win.Height = int32(height)
    })

    return win, nil
}
