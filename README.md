# goal

GoAl -- Go Aliases

Start using:

Install via `brew`:

```shell
brew install #TODO
```

# Idea behind

1. Local alias management  
   To avoid typing repeatable commands
2. AssD - Aliases as a Documentation :D  
   No need to read through whole README file to start operating on you infrastructure

# Project plan

- [ ] Add manual approve step
- [ ] Add "environment" management to avoid tf-plan-dev, tf-plan-stage, tf-plan-prod, etc.
- [ ] Add "depends on" other task like switch to dev?
- [ ] Recursive dependencies
- [ ] Assertions
   - [ ] ref output
   - [ ] recursive assertions 
   - [ ] raw CLI output -- bad pattern?
- [ ] Global aliases in `$HOME` directory?
- [ ] Self-autocompletion via [https://github.com/posener/complete](complete) library
- [ ] Generate ops-doc from commands
- [ ] Support both goal.yaml & goal.yml
- [ ] Support `-f my-goal.yaml`
