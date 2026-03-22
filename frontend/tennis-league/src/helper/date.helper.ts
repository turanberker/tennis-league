export const formatDate = (value: Date | null | undefined): string => {
    if (!value) return "";

    // Eğer string gelirse Date objesine çevir
    const date = value instanceof Date ? value : new Date(value);

    // Geçersiz tarih kontrolü
    if (isNaN(date.getTime())) return "";

    return date.toLocaleDateString('tr-TR', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric'
    });
};