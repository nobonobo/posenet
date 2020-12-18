package posenet

//go:generate sh -c "cd src && yarn install && esbuild index.js --format=esm --global-name=posenet --bundle --platform=node --outfile=../dist/posenet.js"

import (
	"fmt"
	"syscall/js"

	"github.com/nobonobo/spago"
	"github.com/nobonobo/spago/jsutil"
)

var (
	document  = js.Global().Get("document")
	navigator = js.Global().Get("navigator")
	posenet   = spago.LoadModuleAs("posenet", "https://nobonobo.github.io/posenet/dist/posenet.js")
)

// PoseNet ...
type PoseNet struct {
	net    js.Value
	video  js.Value
	Config Config
}

// Config ...
type Config struct {
	Algorithm       string  `js:"algorithm"`
	Architecture    string  `js:"architecture"`
	OutputStride    int     `js:"outputStride"`
	InputResolution int     `js:"inputResolution"`
	Multiplier      float64 `js:"multiplier"`
	QuantBytes      int     `js:"quantBytes"`
}

// DefaultSingleConfig ...
var DefaultSingleConfig = Config{
	Algorithm:       "single-pose",
	Architecture:    "MobileNetV1",
	OutputStride:    16,
	InputResolution: 200,
	Multiplier:      0.5,
	QuantBytes:      2,
}

// DefaultMultipleConfig ...
var DefaultMultipleConfig = Config{
	Algorithm:       "multiple-pose",
	Architecture:    "MobileNetV1",
	OutputStride:    16,
	InputResolution: 200,
	Multiplier:      0.5,
	QuantBytes:      2,
}

// New ...
func New(config Config) *PoseNet {
	return &PoseNet{
		net:    js.Null(),
		video:  js.Null(),
		Config: config,
	}
}

// Start ...
func (n *PoseNet) Start(videoID string) error {
	n.Stop()
	net, err := jsutil.Await(posenet.Call("load"))
	if err != nil {
		return err
	}
	n.net = net
	config := map[string]interface{}{
		"audio": false,
		"video": map[string]interface{}{
			"facingMode": "true",
		},
	}
	stream, err := jsutil.Await(navigator.Get("mediaDevices").Call("getUserMedia", config))
	if err != nil {
		return err
	}
	video := document.Call("getElementById", videoID)
	if video.IsUndefined() {
		return fmt.Errorf("not found video id=%q", videoID)
	}
	video.Set("srcObject", stream)
	success, err := jsutil.Await(
		js.Global().Get("Promise").New(jsutil.Callback1(func(resolve js.Value) interface{} {
			video.Set("onloadedmetadata", jsutil.Callback0(func() interface{} {
				resolve.Invoke(true)
				return nil
			}))
			return nil
		})),
	)
	if err != nil {
		return err
	}
	video.Call("play")
	ok, err := jsutil.Await(
		js.Global().Get("Promise").New(jsutil.Callback1(func(resolve js.Value) interface{} {
			video.Set("onloadeddata", jsutil.Callback0(func() interface{} {
				resolve.Invoke(true)
				return nil
			}))
			return nil
		})),
	)
	if err != nil {
		return err
	}
	if !success.Bool() || !ok.Bool() {
		return fmt.Errorf("video can't ready video stream")
	}
	n.video = video
	return nil
}

// Stop ...
func (n *PoseNet) Stop() {
	if !n.net.IsNull() {
		n.net.Call("dispose")
		n.net = js.Null()
	}
	if !n.video.IsNull() {
		n.video.Get("srcObject").Call("getTracks").Call("forEach", jsutil.Callback1(func(t js.Value) interface{} {
			t.Call("stop")
			return nil
		}))
		n.video = js.Null()
	}
}

// GetAdjacentKeyPoints ...
func GetAdjacentKeyPoints(keypoints js.Value, minConfidence float64) js.Value {
	return posenet.Call("getAdjacentKeyPoints", keypoints, minConfidence)
}

// EstimateSinglePose ...
func (n *PoseNet) EstimateSinglePose(option map[string]interface{}) js.Value {
	if option == nil {
		return n.net.Call("estimateSinglePose", n.video)
	}
	return n.net.Call("estimateSinglePose", n.video, option)
}

// EstimateMultiplePoses ...
func (n *PoseNet) EstimateMultiplePoses(option map[string]interface{}) js.Value {
	if option == nil {
		return n.net.Call("estimateMultiplePoses", n.video)
	}
	return n.net.Call("estimateMultiplePoses", n.video, option)
}
