export async function getResults(searchTerm: string) {
    try {
        console.log("Search term: ", searchTerm);
        const response = await fetch("http://localhost:8080/search", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                SearchTerm: searchTerm,
            }),
        });
        if (!response.ok) {
            console.log("Error: ", response.statusText);
        }

        return await response.json();
    } catch (error) {
        console.log("Fetch error: ", error);
        return [];
    }
}