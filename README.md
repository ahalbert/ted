# fsaed: Finite State Autonoma EDitor
A tool for editing files according to the rules of a provided Finite State Machine 


## Examples

### Run sed only after seeing multiple patterns

Given the input:

```
baz
foo
baz
bar
baz
```

And you only want to edit the final `baz` into `bang`, use this command:

```
$ echo "baz\nfoo\nbaz\nbar\nbaz" | fsaed '/foo/ /bar/ do s/baz/bang/'
```

Results in: 

```
baz
foo
baz
bar
bang
```

### Print Lines Between /regex/

Given the input:

```
DO NOT PRINT THIS LINE
baz - DO NOT PRINT THIS EITHER
foo
bar
baz - DO NOT PRINT THIS EITHER
DO NOT PRINT THIS LINE
```

And you only want to print what's between the `baz`s

```
$ fsaed -n '/baz/ /baz/;print' < file.txt
```

Results in:
```
foo
bar
```

### Run sed if /regexs/ are seen, but reset if /badregex/ is seen

Given the input:

```
beep
boop
buzz
cheater
beep
boop
cheater
```

And you want to modify `cheater` to `nose` only if you see a beep and buzz, but if there's a `buzz`, start looking for `/beep/` again

```
$ fsaed '/beep/ /boop/;/buzz/ -> 1 do s/cheater/nose/ ;/buzz/ -> 1' < file.txt
```

Results In:

```
beep
boop
buzz
cheater
beep
boop
nose
```

## Syntax

fsaed consists of *states*, which contain *actions*. During each execution, `fsaed` will:

1. Read a line from the input.
2. Execute each action for that state in the order parsed
3. If an action requires it to move state, stops executing actions and moves to the next line


### Statement

```
[<statename>:] Action [; Action]
```

Binds the Action to the state `statename`. If a state is not specified, it is an *Anonymous State*, and assigned a name from 1..N, incrementing each time a new state is created. Multiple actions in a statement can be combined using `;`.


### Action

Various actions can be specified in a state:

#### Goto on /regex/

`/<regex>/ [-> <statename>]`

Change current state to state `statename` if input line matches `regex`. If a state is not specified, assigns it to the state `highestAnonymousState + 1`

#### Do Sed Action

`do s/sed/command/g`

Execute sed command on input line.

#### Goto Action

`-> statename`

Change current state to `statename`
