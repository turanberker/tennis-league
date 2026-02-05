import React from 'react';
import { PrimeReactProvider } from 'primereact/api';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import reportWebVitals from './reportWebVitals';
import 'primereact/resources/themes/lara-light-blue/theme.css'; // tema
import 'primereact/resources/primereact.min.css'; // core css
import 'primeicons/primeicons.css'; // prime icons
import '@fortawesome/fontawesome-free/css/all.min.css'; // font awesome

ReactDOM.createRoot(document.getElementById('root')).render(
  <PrimeReactProvider>
    <App />
  </PrimeReactProvider>,
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
