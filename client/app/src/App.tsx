import React from 'react';
import './App.css';
import SearchInput from './components/Search';

function App() {
    const [searchTerm, setSearchTerm] = React.useState('');

    return (
      <html lang="english">
        <body>
          <h1>Welcome to my search engine!</h1>
            <h2>Search for a term</h2>
                <SearchInput
                    searchTerm={searchTerm}
                    setSearchTerm={setSearchTerm}
                />
            {/*<form method="get" action="/search">*/}
            {/*    <label>*/}
            {/*      Enter a search term:*/}
            {/*      <input type="text" name="term"/>*/}
            {/*    </label>*/}
            {/*    <label>*/}
            {/*      <input type="submit" name="submit"/>*/}
            {/*    </label>*/}
            {/*</form>*/}
          </body>
      </html>
  );
}

export default App;
