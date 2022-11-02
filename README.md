## AIOps is a platform consists of micro services and web frontend and all AI(artifitial intelligence) workflow can be done here

It consists of 4 parts, and each parts can work independently:
- dataset center: user-friendly training dataset management, team's dataset can be stored here
- model center: all popular kind of models(PyTorch, Tensorflow, ONNX, etc.) can be managed here and supply transformation between each other.
- task center: we divide task into 2 types: interact task and automatic task. we implement interact task by using jupity notebook
  task center can be used to run model training tasks, data processing tasks and some experiments
- service center: supply model inference or serving

### AIOps is based on k8s to schedule computing resources and juicefs to spare storage resources
