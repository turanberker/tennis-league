import React, { ReactNode } from 'react';
import { useFormContext, get } from 'react-hook-form';
import { classNames } from 'primereact/utils';

interface FormItemProps {
  label?: string;
  name: string; // errors objesindeki anahtar ile eşleşmeli
  children: ReactNode;
  required?: boolean;
  className?: string;
}

const FormItem: React.FC<FormItemProps> = ({
  label,
  name,
  children,
  required = false,
  className = "",
}) => {
  // Formun genelindeki errors ve diğer bilgileri context'ten alıyoruz
  const { formState: { errors } } = useFormContext();

  // Nested (iç içe) objeler için de hata bulmayı kolaylaştıran 'get' fonksiyonu
  const error = get(errors, name);

  return (
    <div className={classNames("flex flex-column gap-2", className)}>
      {label && (
        <label htmlFor={name} className="font-medium text-900">
          {label} {required && <span className="text-red-500">*</span>}
        </label>
      )}

      {children}

      {error && (
        <small className="p-error text-sm">
          {error.message}
        </small>
      )}
    </div>
  );
};

export default FormItem;