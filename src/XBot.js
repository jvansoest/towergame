import React, { useRef, useState, useEffect } from "react";
import { useLoader, useFrame } from "react-three-fiber";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader";
import model3d from "./models/Xbot.glb";
import * as THREE from "three";
import Box from "./Box.js";
//HOW TO USE: <Suspense fallback={<Box />}>{<XBot />}</Suspense>
const XBot = (props) => {
  const mesh = useRef();
  const gltf = useLoader(GLTFLoader, model3d);

  const model = gltf.scene;
  const animations = gltf.animations;

  var clock = new THREE.Clock();
  let delta = 0;
  const mixer = new THREE.AnimationMixer(model);

  useFrame((state, time) => {
    if (!mixer) return;
    //if (!state.animate) return;
    delta = clock.getDelta();
    mixer.update(delta);
  });

  const clip = animations[6];
  const action = mixer.clipAction(clip);
  action.play();

  var vec = new THREE.Vector3();
  vec.set(0, 1, 0);

  model.setRotationFromAxisAngle(vec, Math.PI / 2);

  return <primitive ref={mesh} object={model} position={props.position} />;
  //return <Box />;
};

export default XBot;
