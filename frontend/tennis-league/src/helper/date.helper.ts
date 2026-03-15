export const formatDate = (value: Date | null | undefined): string => {
    if (value) {
        return value.toLocaleDateString('tr-TR', {
            day: '2-digit',
            month: '2-digit',
            year: 'numeric'
        });
    }
    return ""; // Explicitly return the empty string
};