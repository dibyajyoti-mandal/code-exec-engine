import React, { useRef, useState } from "react";
import Editor from "@monaco-editor/react";

const defaultCode = {
  javascript: `// JavaScript demo
console.log('Hello from Monaco!');`,

  typescript: `// TypeScript demo
const greet = (name) => "Hello " + name;
console.log(greet("World"));`,

  python: `# Python demo (cannot run in browser)
print("Hello from Python")`,

  cpp: `#include <iostream>
using namespace std;
int main() {
    cout << "Hello C++";
    return 0;
}`,

  java: `public class Main {
    public static void main(String[] args) {
        System.out.println("Hello Java");
    }
}`,
};

export default function App() {
  const [language, setLanguage] = useState("javascript");
  const [code, setCode] = useState(defaultCode["javascript"]);
  const [output, setOutput] = useState("");
  const editorRef = useRef(null);

  const handleEditorDidMount = (editor) => {
    editorRef.current = editor;
  };

  const handleLanguageChange = (e) => {
    const lang = e.target.value;
    setLanguage(lang);
    setCode(defaultCode[lang] || "");
  };

  const runInBrowser = () => {
    if (language !== "javascript") {
      setOutput("Browser-run demo works only for JavaScript.");
      return;
    }

    const codeValue = editorRef.current?.getValue() || code;

    const logs = [];
    const originalLog = console.log;

    console.log = (...args) => logs.push(args.join(" "));

    try {
      eval(codeValue);
    } catch (err) {
      logs.push("Error: " + err.toString());
    }

    console.log = originalLog;
    setOutput(logs.join("\n"));
  };

  const submitToServer = async () => {
    const payload = {
      language,
      code: editorRef.current?.getValue() || code,
      problemId: "demo-1",
    };

    setOutput("Submitting...");

    try {
      const res = await fetch("/api/submit", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      const data = await res.json();
      setOutput("Server response: \n" + JSON.stringify(data, null, 2));
    } catch (err) {
      setOutput("Error submitting: " + err.toString());
    }
  };

  return (
    <div style={{ height: "100vh", display: "flex", flexDirection: "column" }}>
      <div
        style={{
          padding: "12px",
          background: "#1e1e1e",
          color: "white",
          display: "flex",
          justifyContent: "space-between",
        }}
      >
        <h2>Code Editor</h2>

        <div style={{ display: "flex", gap: "10px" }}>
          <select value={language} onChange={handleLanguageChange}>
            <option value="javascript">JavaScript</option>
            <option value="typescript">TypeScript</option>
            <option value="python">Python</option>
            <option value="cpp">C++</option>
            <option value="java">Java</option>
          </select>

          <button onClick={runInBrowser}>Run Code</button>
          <button onClick={submitToServer}>Submit Code</button>
        </div>
      </div>

      <div style={{ display: "flex", flexGrow: 1 }}>
        <div style={{ width: "50%", borderRight: "1px solid #ddd" }}>
          <Editor
            height="100%"
            language={language}
            value={code}
            theme="vs-dark"
            onChange={(val) => setCode(val)}
            onMount={(editor) => handleEditorDidMount(editor)}
            options={{ fontSize: 14 }}
          />
        </div>

        {/* Output */}
        <div style={{ width: "50%", padding: "16px" }}>
          <h3>Output</h3>
          <pre
            style={{
              background: "#111",
              color: "#0f0",
              padding: "12px",
              borderRadius: "6px",
              height: "80%",
              overflow: "auto",
            }}
          >
            {output}
          </pre>
        </div>
      </div>
    </div>
  );
}
