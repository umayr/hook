# Hook
> A Github webhook handler primarily to notify about your PRs.

## Preface
Our team relies on Github PRs and reviews to release new stuff and bug fixes as many teams do; major bottleneck in this scenario is that people don't review quite often even if PRs are assigned to them and they need to nudge about pending reviews. I, personally wouldn't want to be a bottleneck for anybody's dear feature. Github does notify you via email but it usually gets piled up. To resolve this, I initially tried to allocate two slots in my calendar for reviews, it doesn't really work out. So I'm trying something dynamic which works like this:

- A new review is assigned to me, I get a notification.
- A new slot for that particular review would be added to my calendar (based on the periods I have provided)
- No review pending, I wouldn't have anything in my calendar, I would do whatever I have been doing.
- Any review pending, at least I'd get notification and time allocated.

I just started doing this, so I'm not sure it would make things better for me. This repository has the first part implemented for now, as I'm still testing out the calendar stuff. 

## Setup
Assuming you've Go set up in your machine:
```bash
  # clone the repository
  git clone git@github.com:umayr/hook.git
  
  # make the binary
  make
  
  # create a zip file
  make zip
```
## Deployment
I'm using aws lambda for this, you can use whatever you want but you would have to write server yourself (don't forget to send a PR for that). For lambda you can do something like this:
```bash
  aws lambda create-function \
    --region region \
    --function-name lambda-handler \
    --memory 128 \
    --role arn:aws:iam::account-id:role/execution_role \
    --runtime go1.x \
    --zip-file fileb://path-to-your-zip-file/handler.zip \
    --handler lambda-handler
```

More information could be found [here](https://docs.aws.amazon.com/lambda/latest/dg/lambda-go-how-to-create-deployment-package.html)

## Notes
This only has implementation to send notification via pushover, if you want to use something else, you can simply create a new notifier, like:
```go
type Whatever struct {}

func (w *Whatever) Notify(title, message string) error {
    // sends notification
    return nil
}

func main() {
    // parse the payload
    p := Parse(payload)
    h := hook.NewHook(p, &Whatever{})
    
    err := h.Perform()
    // handle error
}
```
