# TODO

- [ ] Create a file server that will run as a separate process in the worker's environment
- [ ] Dockerise the worker process, the container should start both the worker and the file server
- [ ] Update the worker to communicate the endpoint URL of the file server serving the downloaded files to the core web server
- [ ] Add endpoint in the core web server to poll the status of the checkout job. Upon completion, it'll read the file server's download URL from some destination, and it'll send that to the client
- [ ] At this point, we can start creating the frontend. Ultimately, once we get around to connecting the aforementioned last few pieces, the frontend will download the audio files from the server sent file server URL
