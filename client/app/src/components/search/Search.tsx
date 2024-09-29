import React from "react";
import SearchInput from "./SearchInput";
import Button from '@mui/material/Button';
import {getResults} from "../../api/search";

interface SearchComponentProps {
    searchTerm: string;
    setSearchTerm: React.Dispatch<React.SetStateAction<string>>;
    setResults: React.Dispatch<React.SetStateAction<any[]>>;
}

export default function Search(props: SearchComponentProps) {
    const handleSearch = async () => {
        const response = await getResults(props.searchTerm);
        props.setResults(response.HITS);
    }

    return (
        <div className="flex flex-row gap gap-x-6">
            <SearchInput
                searchTerm={props.searchTerm}
                setSearchTerm={props.setSearchTerm}
            />
            <Button variant="outlined" onClick={handleSearch}>
                Search
            </Button>
        </div>
    )
}