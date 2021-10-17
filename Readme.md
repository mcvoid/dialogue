# Dialogue Engine

## Script Format

A script is a markdown-like collection of dialogue nodes. These nodes are the different sequences of
dialogue to be displayed to the user. A node looks like this:

````markdown
// the first section is the frontmatter. This contains type definitions for external functions.
// These are defined as follows:
extern func1(bool, number, string, null);
extern func2();

// end frontmatter with three backticks.
```

# node1

Here's some text that's going to be displayed to the user.
Just like in HTML, extra whitespace including newlines are
reduced down to render as a single space.

Blocks of text separated by an empty line are displayed to
the user as separate chunks of text, so the user will see
the above paragraph, then this one.

[node2](Markdown links instruct the engine to transition to another node)

# node2

Nodes can be transitioned by links and by giving an option.
Links don't give an option - they just display some text and
the transition.

The following shows what an option looks like. In markdown terms,
it is just an unordered list of links.

- [node1] (This goes to node1)
- [node2] (This repeats node2)
- [node3] (An empty line will terminate the list, as it does paragraphs and links)

# node3

Paragraph text can have inline code interleaved with the plain text.
This is useful to print out the value of a variable like so: `variable1`.
Just use the markdown inline code element. It can also have simple expressions
like so: `3 + 3`

```
// You can also have code blocks.
// These are Markdown fenced code blocks.
// They are a block-level element like a paragraph, link, or option list.
// The language in the code blocks is a simple C-like language.
// Comments are like C single-line comments
// You can assign variables

variable1 = 32 * variable2;

// If statements  and if-else statements are in there as well
if variable1 > 12 {
  // You can call functions that the vm exposes
  function1(variable1 == variable2);
}

// There's loops
while variable1 > 12 {
  // preincrement and predecrement are expressions and don't reassign the referenced variable
  variable1 = --variable1
}

// you can transition to other nodes
goto node4;

// you can call the external functions
func1(false, 5, "abc", null);

```

# node4

The first node in the script is the dialogue's starting point.
Ending a node without a link or choice will terminate the script.

````

## Todo List

- [x] unit test the parser
- [x] unit test the parse tree -> ast translation
- [x] unit test type checking
- [x] implement constant folding
- [x] unit test constant folding
- [x] implement test dead node elimination
- [x] unit test dead node elimination
- [x] 100% unit test coverage everywhere (internal + user-facing)
- [x] add if with elided else code generation
- [x] add naked if with reverse logic when consequent is empty
- [x] add unary minus operator
- [x] add != operator
- [x] add modulo operator
- [x] make infinite loop AST node when while(true){...}
- [x] add extern func declarations
- [x] add optional front matter code block
- [x] implement parent package/public interface
- [x] implement compiler CLI
- [x] implement vm cli
- [ ] make link text support inline expressions
- [ ] implement asm deflate/inflate for binary format
- [ ] (stretch goal) Multi-file script linker
- [ ] (stretch goal) break and continue
- [ ] (stretch goal) variable scopes / differentiating extern and in-script variables
- [ ] (stretch goal) add typed variable declarations
- [ ] (stretch goal) add builtin function calls (EndDialog, PushOption, ShowOption, ShowLine, EnterNode, ExitNode)
