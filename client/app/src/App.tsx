import React, {useEffect} from 'react';
import './App.css';
import Search from "./components/search/Search"
import { getResults } from "./api/search";

function App() {
    const [searchTerm, setSearchTerm] = React.useState('');
    const [results, setResults] = React.useState<any[]>([]);

    useEffect(() => {
        console.log("Search Term: ", searchTerm);
        getResults(searchTerm).then((response: any) => {
            console.log("Response: ", response);
            setResults(response);
        });

    }, [searchTerm]);


    return (
      <html lang="english">
        <body>
          <h1>Welcome to my search engine!</h1>
            <h2>Search for a term</h2>
                <Search
                    searchTerm={searchTerm}
                    setSearchTerm={setSearchTerm}
                />
          {results.length > 0 && (
              <div>
                  {results.map((item: any) => (item))}
              </div>
          )}
          </body>
      </html>
  );
}

export default App;
