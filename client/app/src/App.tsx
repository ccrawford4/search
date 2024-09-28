import React from 'react';
import './App.css';
import Search from "./components/search/Search"

function App() {
    const [searchTerm, setSearchTerm] = React.useState('');

    return (
      <html lang="english">
        <body>
          <h1>Welcome to my search engine!</h1>
            <h2>Search for a term</h2>
                <Search
                    searchTerm={searchTerm}
                    setSearchTerm={setSearchTerm}
                />
          </body>
      </html>
  );
}

export default App;
