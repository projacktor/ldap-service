'use client'

import ParticleBackground from 'react-particle-backgrounds'

const particleSettings = {
  canvas: {
    canvasFillSpace: true,
    width: 200,
    height: 200,
    useBouncyWalls: false
  },
  particle: {
    particleCount: 150,
    color: '#fff',
    minSize: 1,
    maxSize: 5
  },
  velocity: {
    minSpeed: 0.2,
    maxSpeed: 0.5
  },
  opacity: {
    minOpacity: 0,
    maxOpacity: 0.6,
    opacityTransitionTime: 10000
  }
}

export default function BackgroundWrapper() {
  return (
    <div className="absolute inset-0 -z-10">
      <ParticleBackground settings={particleSettings} />
    </div>
  )
}
