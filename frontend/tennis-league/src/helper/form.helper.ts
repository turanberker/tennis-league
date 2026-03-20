
// Şemayı ve alan adını verince true/false döner
export const isFieldRequired = (schema: any, fieldName: string) => {
    return schema
        .describe()
        .fields[fieldName]?.tests.some((test: any) => test.name === "required");
};