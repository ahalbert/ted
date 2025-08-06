import {EditorView, basicSetup} from "codemirror"
import { EditorState } from '@codemirror/state';
import {parser} from "./ted_parser.js"
import {foldNodeProp, foldInside, indentNodeProp} from "@codemirror/language"
import {styleTags, tags as t} from "@lezer/highlight"
import {LRLanguage} from "@codemirror/language"
import {completeFromList} from "@codemirror/autocomplete"
import {LanguageSupport} from "@codemirror/language"

export const parserWithMetadata = parser.configure({
  props: [
    styleTags({
      "start stop capture print println do dountil clear rewind fastforward if else" : t.keyword,
      identifier: t.tagName,
      stateidentifier: t.variableName,
      String: t.string,
      Regex: t.regexp,
      Boolean: t.bool,
      String: t.string,
      LineComment: t.lineComment,
      "( )": t.paren
    }),
    indentNodeProp.add({
      Application: context => context.column(context.node.from) + context.unit
    }),
    foldNodeProp.add({
      Application: foldInside
    })
  ]
});



export const exampleLanguage = LRLanguage.define({
  name: "ted",
  parser: parserWithMetadata,
  languageData: {
        commentTokens: {line: "#"}
  }
});

export function ted() {
    return new LanguageSupport(exampleLanguage,  []);
}

function createEditorStateForTed(initialContents, options = {}) {
    let extensions = [
      basicSetup,
      ted()
    ];

    return EditorState.create({
        doc: initialContents,
        extensions
    });
}

function createEditorState(initialContents, options = {}) {
    let extensions = [
      basicSetup
    ];

    return EditorState.create({
        doc: initialContents,
        extensions
    });
}

function createEditorView(state, parent) {
    return new EditorView({ state, parent });
}

export { createEditorStateForTed, createEditorState, createEditorView};


