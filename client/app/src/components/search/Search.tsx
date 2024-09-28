import React from "react";
import SearchInput from "./SearchInput";
import Button from '@mui/material/Button';

interface SearchComponentProps {
    searchTerm: string;
    setSearchTerm: React.Dispatch<React.SetStateAction<string>>;
}

export default function Search(props: SearchComponentProps) {
    return (
        <div className="flex flex-row gap gap-x-6">
            <SearchInput
                searchTerm={props.searchTerm}
                setSearchTerm={props.setSearchTerm}
            />
            <Button
                variant="outlined"
                onClick={() => props.setSearchTerm(props.searchTerm)}
            >
                Search
            </Button>
        </div>
    )
}