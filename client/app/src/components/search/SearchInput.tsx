import * as React from 'react';
import TextField from '@mui/material/TextField';
import Autocomplete from '@mui/material/Autocomplete';

interface SearchInputProps {
    searchTerm: string;
    setSearchTerm: React.Dispatch<React.SetStateAction<string>>;
}

export default function SearchInput(props: SearchInputProps) {
    const emptyParams: any = [props.searchTerm];
    return (
        <Autocomplete
            freeSolo
            disablePortal
            options={emptyParams}
            sx={{ width: 300 }}
            onSelect={(event) => props.setSearchTerm((event.target as HTMLSelectElement).value)}
            renderInput={(params: any) =>
                <TextField
                    {...params}
                    label="Search Term"
                />
            }
        />
    );
}