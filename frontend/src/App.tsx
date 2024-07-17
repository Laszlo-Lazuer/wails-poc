import {ChangeEvent, useEffect, useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Greet} from "../wailsjs/go/main/App";
import { Add } from "../wailsjs/go/main/App";

function App() {
    const[input, setInput] = useState({
        num1: "",
        num2: ""
    })

    const [result, setResult] = useState("");

    useEffect(() => {
        Add(+input.num1, +input.num2).then((v) => setResult(String(v)));
    }, [input])

    function handleChange(event: ChangeEvent<HTMLInputElement>): void {
        setInput({
            ...input,
            [event.target.name]: event.target.value
        })
    }

    return (
        <div>
            <input name="num1" value={input.num1} type="text" onChange={handleChange}></input>
            <input name="num2" value={input.num2} type="text" onChange={handleChange}></input>


            <h1> Result is {result}</h1>
        </div>
    )
}

export default App
