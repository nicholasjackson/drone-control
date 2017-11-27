package control

import (
	"log"
	"time"

	messages "github.com/nicholasjackson/drone-messages"
)

// StartFollowing allows the drone to start tracking a face and follow it
func (a *AutoPilot) StartFollowing() {
	log.Println("Start Following")
	a.following = true
	a.lastFace = nil
}

// StopFollowing stops tracking a face
func (a *AutoPilot) StopFollowing() {
	log.Println("Stop Following")
	a.following = false
	a.lastFace = nil
	a.drone.Stop()
}

// FollowFace moves to drone in the direction of a detected face
func (a *AutoPilot) FollowFace(m *messages.FaceDetected) {
	if a.lastFace != nil && a.following {
		a.moveDrone(m)
	}

	a.lastFace = m
}

func (a *AutoPilot) moveDrone(m *messages.FaceDetected) {
	log.Println("Got Face, moving...", m)

	// calculate the right moves
	centerPoint := (m.Bounds.Max.X) / 2
	faceCenter := ((m.Faces[0].Max.X - m.Faces[0].Min.X) / 2) + m.Faces[0].Min.X

	log.Println("Centre:", centerPoint)
	log.Println("Face Center:", faceCenter)

	if faceCenter < (centerPoint - a.minDistance) {
		log.Println("Left")
		if !a.movingLeft {
			a.drone.Left(a.speed)
			a.movingLeft = true
		}
	} else if (centerPoint + a.minDistance) < faceCenter {
		log.Println("Right")
		if !a.movingRight {
			a.drone.Right(a.speed)
			a.movingRight = true
		}
	} else {
		log.Println("Stop")
		a.movingLeft = false
		a.movingRight = false
		a.drone.Stop()
	}

	a.setDeadMansSwitch()
}

// setDeadMansSwitch sets a stop command after timeout incase no further face
// tracking info is received
func (a *AutoPilot) setDeadMansSwitch() {
	log.Println("DMS Set")
	if a.deadMansSwitch == nil {
		log.Println("Start DMS", a.timeout)

		a.deadMansSwitch = time.AfterFunc(a.timeout, func() {
			a.drone.Stop()
		})

		return
	}

	if !a.deadMansSwitch.Stop() {
		<-a.deadMansSwitch.C
	}
	a.deadMansSwitch.Reset(a.timeout)
}
