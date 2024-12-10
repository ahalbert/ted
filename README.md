# fsaed: Finite State Autonoma EDitor
A tool for editing files according to the rules of a provided Finite State Machine 

## Motivation

Once, I was presented with an the following file (example abridged) 

```
INFO:2024-12-07 13:01:40:Trace:198d079c-af9a-45b2-8236-7fbb2a012f69:Starting Procedure foo
ERROR:2024-12-07 13:01:41:Trace:198d079c-af9a-45b2-8236-7fbb2a012f69:Error 1
INFO:2024-12-07 13:01:41:Trace:198d079c-af9a-45b2-8236-7fbb2a012f69:Ending Procedure foo
INFO:2024-12-07 13:01:41:Trace:198d079c-af9a-45b2-8236-7fbb2a012f69:Starting Procedure bar
INFO:2024-12-07 13:01:41:Trace:198d079c-af9a-45b2-8236-7fbb2a012f69:Error 2
INFO:2024-12-07 13:01:41:Trace:198d079c-af9a-45b2-8236-7fbb2a012f69:Success
INFO:2024-12-07 13:01:42:Trace:198d079c-af9a-45b2-8236-7fbb2a012f69:Ending Procedure bar
INFO:2024-12-07 13:01:42:Trace:30019fff-7645-4d07-9fc4-0bbb39aa09db:Starting Procedure foo
INFO:2024-12-07 13:01:42:Trace:30019fff-7645-4d07-9fc4-0bbb39aa09db:Success
INFO:2024-12-07 13:01:42:Trace:30019fff-7645-4d07-9fc4-0bbb39aa09db:Ending Procedure foo
INFO:2024-12-07 13:01:43:Trace:30019fff-7645-4d07-9fc4-0bbb39aa09db:Starting Procedure bar
ERROR:2024-12-07 13:01:43:Trace:30019fff-7645-4d07-9fc4-0bbb39aa09db:Error 3
ERROR:2024-12-07 13:01:43:Trace:30019fff-7645-4d07-9fc4-0bbb39aa09db:Error 4
INFO:2024-12-07 13:01:44:Trace:30019fff-7645-4d07-9fc4-0bbb39aa09db:Ending Procedure bar

```

I wanted only the errors that did not have a success in the procedure. In this case, we should only get Errors 1,3,4

```
ERROR:2024-12-07 13:01:41:Trace:198d079c-af9a-45b2-8236-7fbb2a012f69:Error 1
ERROR:2024-12-07 13:01:43:Trace:30019fff-7645-4d07-9fc4-0bbb39aa09db:Error 3
ERROR:2024-12-07 13:01:43:Trace:30019fff-7645-4d07-9fc4-0bbb39aa09db:Error 4 
```


I created an `awk` program to keep track of things and get the correct output. But I thought it was easier to express what I wanted as a state machine. Thus was born `fsaed`, a language for specifying state machines and using them to process files. 

An equivalent `fsaed` program:

```
startstate: /Starting.Procedure/ -> capture_begin
capture_begin: { 
	start capture 
	-> lookforsuccessorending 
	/Success/ -> startstate
}
lookforsuccessorending: /Success/  -> startstate 
lookforsuccessorending: /Ending.Procedure/ { 
	stop capture 
	print 
	-> startstate 
}

```

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
$ echo "baz\nfoo\nbaz\nbar\nbaz" | fsaed '/foo/ -> /bar/ -> do s/baz/bang/'
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
$ fsaed -n 'stop:/baz/ -> start start:/baz/ -> 1 start: print' < file.txt
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
$ fsaed '/beep/ -> {/boop/ -> /buzz/ -> 1} {do s/cheater/nose/ /buzz/ -> 1}' < file.txt
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

### Capturing 

*Capturing* enables you to read input into a variable rather than printing it on the screen. 

#### Capture single line

Given the input:

```
beep
boop
foo
bar
baz
buzz
```

You can capture one line as so:

```
1: /beep/ ->
2: {capture mycapture -> }
3: do s/THIS.IS.CAPTURED/CAPTURED/ mycapture
3: /buzz/ -> 
4: print mycapture

```

This program removes the `boop`, captured into the variable `$_`:

```
beep
boop
foo
buzz
CAPTURED
bar
CAPTURED
baz
```

#### Capture multiple lines

Given the input:

```
beep
boop - CAPTURED
foo - CAPTURED
bar - CAPTURED
baz
buzz
```

And running this `fsaed` program with `--no-print` option:

```
/beep/ ->
/boop/ {start capture ->} 
/baz/ {stop capture print -> 1}
```


Yields:

```
boop - CAPTURED
foo - CAPTURED
bar - CAPTURED
```

### Regex Capture Groups

You can store capture groups in a variable and refer to them later.

Given the input: 

```
beep
boop
i want these and those
foo
bar
baz
buzz
```

And this program with the `--no-print` option: 

```
/i.*want.(these).and.(those)/ ->
{println $1 println $2 ->}
```

Yields:

```
these
those
```

## Syntax

fsaed consists of *states*, which contain *actions*. During each execution, `fsaed` will:

1. Read a line from the input.
2. Execute each action for that state in the order parsed
3. If an action requires it to move state, stops executing actions and moves to the next line.
4. Prints a line unless `--no-print` or capturing is on.


### Statement

```
[<statename>:] Action [, Action]
```

Binds the Action to the state `statename`. If a state is not specified, it is an *Anonymous State*, and assigned a name from 1..N, incrementing each time a new state is created. Multiple actions in a statement can be combined using `{ }`. If you want to specify multiple different rules for the same state, use `,`


### Action

Various actions can be specified in a state:

#### Do action on Regex

`/<regex>/ Action`

Perform `Action` if the current line matches regex.

#### Do Sed Action

`do s/sed/command/g [variable]`

Execute `sed` command on `variable`. If no `variable` is specified, assumes the current line or capture. 

#### Goto Action

`-> [statename]`

Change current state to `statename`. If a state is not specified, assumes the next state listed in the program. If this is the last state, goes to state "0". 


#### Do multiple actions

`{ Action... }`

Runs all the actions between the `{` and `}`. If state changes, stops executing block.

#### Print

`print [variable]`

Prints `variable`. If a variable is not specified, uses `$_` which can be the current line or capture.


#### Capture

`[start|stop] capture [variable]`

Starts/Stops capturing to `variable`. When capturing is started, input lines are redirected to . If variable is not specified, defaults to `$_`. If `start|stop` is not given, only captures the current line. 

### Predefined Variables

* `$_` The default variable used by arguments. At the beginning of an iteration, stores the current line in `$_` unless it is being used to capture. 
* `$@` Contains the original line read in during the iteration.
* `$0` Contains the matched text of the last regex compared.
* `$1..$N` Contains the first to N capture groups in the last regex compared
