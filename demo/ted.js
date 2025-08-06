const introduction = "# Welcome to ted: Turing EDitor\n# ted allows you to process text files according to rules you specify.\n# Github: https://github.com/ahalbert/ted\n\n1: /foo/ -> 2\n# This is state 1, the starting state. During each cycle of execution, \n# ted reads a line of the input file and excutes the actions of the state \n# it's in. It then prints the line of the input file.\n\n2: do s/buzz/boop/\n# In state 2 we execute a sed command on the input line. buzz is changed to \n# boop. Notice in the output, the first buzz, before foo is still there, while \n# the second buzz is changed to boop - it came after foo.\n";
const introductioninput = "beep\nboop\nbuzz\nfoo\nbeep\nboop\nbuzz";

const begin = "# By default, ted prints each line of the input.\n# You can disable this by setting the $PRINTMODE\n# variable to \"noprint\".\n\nBEGIN: let $PRINTMODE = \"noprint\"\n# Special state run at the beginning of every ted program\n\n1: /foo/ -> 2\n\n2: /bar/ -> 0\n# 0 is a special state, with no directions. \n# 0 cannot be left once entered.\n\n2: println\n# Multiple actions can be specified for each state.\n# Ted executes them in order unless state is changed,\n# at which point it stops the action cycle.";
const anonymous = '# Another way you could construct the previous program\n# is using anonymous states. They are useful in one-liners.\n\nBEGIN: let $PRINTMODE = "noprint"\n\n/foo/ -> \n# If a state is not given, then it is assigned a number\n# starting from 1. This is given state 1. If a state is\n# not specified after -> it goes to the next state listed\n# in the program.\n\n/bar/ ->, println\n# Multiple actions can be specified using a ,\n# If there is no next state in the program, -> transitions\n# to state 0.';
const begininput = "...\nStuff we don't want\n...\nfoo\nbeep\nboop\nbuzz\nbar\n...\nMore stuff we don't want\n..."

const capture = "# Ted can store regex capture groups in a variable for later use.\n\n1: /target:.\"(.*)\"/ {let myvar = $1 -> 2 }\n# {} signifies actions that should be executed together.\n# Even if the state transitions inside the {}, ted\n# will execute all the actions inside it.\n\n\n2: do s/buzz/{{.myvar}}/\n# We stored the first capture group \"(.*)\" in myvar. We then substitute\n# it in the sed expression to change it. \n";
const captureinput = "target: \"foo\"\nbeep:\n  boop: \"buzz\"\nfoo: \n  bar: \"baz\"";

const rewind = '# ted allows you to go backwards and forwards in a file.\n# This program finds the three divs of context around the\n# regular expression /target/\nBEGIN: {let $PRINTMODE = "noprint"}\n\n/target/ { let count = 3 -> }\n\nif count == 0  { \n  let count = 6 #Prep for the next state\n  start capture myvar \n  -> \n} \nelse {\n  rewind /div/ #Move file back to the last div seen\n  let count = count - 1 \n}\n\n/div/ {let count = count - 1}, if count == 0 {stop capture -> } \n# Capture 6 divs worth of context surrounding /target/\n\nEND: { println myvar }\n'
const rewindinput = '<div id="not this one">\n  <div id="this two">\n    <div id="this three">\n      <div id="this four">\n      target\n      </div>\n    </div>\n  </div>\n</div>\n'

const capturing = '# Ted allows you to capture the input read each cycle\n# into a variable. In this example, we are given a log \n# file that has various runs, each given a Trace id. We\n# want the errors \n# errors in each trace that DO NOT have "Success" in them.\n\nBEGIN: {let $PRINTMODE = "noprint"}\n\nstartstate: /Starting.Procedure/ -> capturebegin\n\ncapturebegin: { \n   start capture cap \n   -> lookforending \n}\n# Start capturing input into variable cap. Each line\n# will be read into the variable cap until a "stop capture"\n# action is encountered.\n\nlookforending: /Ending.Procedure/ {\n  print cap\n  stop capture \n  clear cap\n  -> startstate \n}\n# Stops the capture, then empties the variable cap\n\n\nALL: /Success/ { \n   stop capture \n   clear cap\n   -> startstate\n}\n# ALL is a rule that is applied to every state. It is \n# evaluated after the current state\'s actions are applied,\n# and can transition even if one transitioned on the current\n# state\'s actions.'
const capturinginput = 'INFO:2024-12-07 13:01:40:Trace:198d079c-af9a-45b2-8236-7fbb2a012f69:Starting...\nINFO:2024-12-07 13:01:40:Trace:198d079c-af9a-45b2-8236-7fbb2a012f69:Starting Procedure foo\nERROR:2024-12-07 13:01:41:Trace:198d079c-af9a-45b2-8236-7fbb2a012f69:Error 1\nINFO:2024-12-07 13:01:41:Trace:198d079c-af9a-45b2-8236-7fbb2a012f69:Ending Procedure foo\nINFO:2024-12-07 13:01:41:Trace:198d079c-af9a-45b2-8236-7fbb2a012f69:Starting Procedure bar\nINFO:2024-12-07 13:01:41:Trace:198d079c-af9a-45b2-8236-7fbb2a012f69:Error 2\nINFO:2024-12-07 13:01:41:Trace:198d079c-af9a-45b2-8236-7fbb2a012f69:Success\nINFO:2024-12-07 13:01:42:Trace:198d079c-af9a-45b2-8236-7fbb2a012f69:Ending Procedure bar\nINFO:2024-12-07 13:01:42:Trace:30019fff-7645-4d07-9fc4-0bbb39aa09db:Starting...\nINFO:2024-12-07 13:01:42:Trace:30019fff-7645-4d07-9fc4-0bbb39aa09db:Starting Procedure foo\nINFO:2024-12-07 13:01:42:Trace:30019fff-7645-4d07-9fc4-0bbb39aa09db:Success\nINFO:2024-12-07 13:01:42:Trace:30019fff-7645-4d07-9fc4-0bbb39aa09db:Ending Procedure foo\nINFO:2024-12-07 13:01:43:Trace:30019fff-7645-4d07-9fc4-0bbb39aa09db:Starting Procedure bar\nERROR:2024-12-07 13:01:43:Trace:30019fff-7645-4d07-9fc4-0bbb39aa09db:Error 3\nERROR:2024-12-07 13:01:43:Trace:30019fff-7645-4d07-9fc4-0bbb39aa09db:Error 4\nINFO:2024-12-07 13:01:44:Trace:30019fff-7645-4d07-9fc4-0bbb39aa09db:Ending Procedure bar\n  '

const dountil = '# You may want to perform an action only when a succesful\n# substitution is performed. The dountil action allows you \n# to do that. \n\ndountil s/foo/bang/ -> \n\n# Only substitute the first foo seen, then transition to\n# state 0.'
const dountilinput = 'beep:\n   boop:\n      buzz: foo\nfoo:\n   bar:\n      baz: foo';

const fsastate = cm6.createEditorStateForTed(introduction);
const fsaed = cm6.createEditorView(fsastate, document.getElementById("fsaeditor"));
const inputstate = cm6.createEditorState(introductioninput);
const inputed = cm6.createEditorView(inputstate, document.getElementById("inputeditor"));
const outputstate = cm6.createEditorState("");
const outputed = cm6.createEditorView(outputstate, document.getElementById("outputeditor"));

function runProgram() {
  var body = {
    "program" : fsaed.state.doc.toString(),
    "data" : inputed.state.doc.toString()
  };
  $.ajax("https://us-east1-ahalbert-clickstream.cloudfunctions.net/ted-api", {
      method : 'POST',
      data : JSON.stringify(body),
      contentType : 'application/json',
      success: function(data) {
        var newState = cm6.createEditorState(data['output'])
        outputed.setState(newState);
      }
    }
  );
}

function changeExample(fsa, inp) {
  var newState = cm6.createEditorStateForTed(fsa);
  fsaed.setState(newState);
  newState = cm6.createEditorState(inp);
  inputed.setState(newState);
}
