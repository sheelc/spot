# Spot: An Alfred Workflow for Spotify

A new Spotify workflow for Alfred. Still in initial construction but aims to prioritize speed (local functions over API calls) to make the workflow feel snappier.

## Setup Instructions
1. Install the workflow
1. Setup an app as a Spotify developer [here](https://developer.spotify.com/my-applications/#!/applications). Add http://localhost:11075/callback as a redirect URI. Save the changes and keep this page open as you will need the Client Id and the Client Secret on the next step.
1. Type in the keyword 'spot' to launch the workflow in Alfred. It will first ask for the Client Id and then the Client Secret. Finish entering both of those.
1. Enter the keyword 'spot' and press enter on authenticating. This will launch your browser to request permissions from Spotify.

That's it!

#### Acknowledgements
---
This workflow is mainly based on the work in [Spotify for Alfred](https://github.com/citelao/Spotify-for-Alfred) written by [citelao](https://github.com/citelao) and uses the same interface/Spotify tricks as in that workflow. This is mainly an exercise in determining if we can make it faster!
