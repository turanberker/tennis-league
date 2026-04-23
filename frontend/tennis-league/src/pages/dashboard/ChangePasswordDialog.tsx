import React, { useState } from "react";
import { Dialog } from "primereact/dialog";
import { Button } from "primereact/button";
import { Password } from "primereact/password"; // Daha güvenli bir giriş için
import { classNames } from "primereact/utils";
import * as Yup from 'yup';
import { Controller, FormProvider, useForm } from "react-hook-form";
import { ChangePasswordRequest } from "../../model/user.model";
import { yupResolver } from "@hookform/resolvers/yup";
import { changeMyPassword, } from "../../api/profileService";
import { useToast } from "../../hooks/useToast";
import { isFieldRequired } from "../../helper/form.helper";

import FormItem from "../../components/FormItem";

interface ChangePasswordDialogProps {
    visible: boolean;
    onHide: () => void;
}

const changePasswordSchema = Yup.object().shape({
    currentPassword: Yup.string()
        .required('Mevcut şifre gereklidir'),
    newPassword: Yup.string()
        .min(8, 'Şifre en az 8 karakter olmalıdır')
        .required('Yeni şifre gereklidir'),
    confirmPassword: Yup.string()
        .oneOf([Yup.ref('newPassword')], 'Şifreler eşleşmiyor') // KRİTİK NOKTA
        .required('Şifre onayı gereklidir')
});

export default function ChangePasswordDialog({ visible, onHide }: ChangePasswordDialogProps) {

    const { show } = useToast();
    const methods = useForm<ChangePasswordRequest>({
        resolver: yupResolver(changePasswordSchema),
        defaultValues: {
            currentPassword: "",
            newPassword: "",
            confirmPassword: ""
        },
    });

    const onSubmit = async (data: ChangePasswordRequest) => {
        const message = await changeMyPassword(data);
        if (message) {
            show({
                severity: "success",
                summary: "Başarılı",
                detail: "Lig başarıyla oluşturuldu",
            })
            methods.reset();
            onHide();
        };

    }


    return (
        <Dialog
            header="Şifreyi Güncelle"
            visible={visible}
            style={{ width: '25vw' }}
            breakpoints={{ '960px': '75vw', '641px': '100vw' }}
            onHide={() => { methods.reset(); onHide(); }} // Kapatınca formu temizle
            draggable={false}
            resizable={false}
        >
            <FormProvider {...methods}>
                <form onSubmit={methods.handleSubmit(onSubmit)} className="p-fluid">

                    <FormItem
                        label="Mevcut Şifre"
                        name="currentPassword"
                        required={isFieldRequired(changePasswordSchema, "currentPassword")}
                    >
                        <Controller
                            name="currentPassword"
                            control={methods.control}
                            render={({ field, fieldState }) => (
                                <Password
                                    id={field.name}
                                    {...field}
                                    feedback={false}
                                    toggleMask
                                    autoFocus
                                    className={classNames({ "p-invalid": fieldState.error })}
                                    inputClassName="w-full"
                                />
                            )}
                        />
                    </FormItem>
                    <FormItem
                        label="Yeni Şifre"
                        name="newPassword"
                        required={isFieldRequired(changePasswordSchema, "newPassword")}
                    >
                        <Controller
                            name="newPassword"
                            control={methods.control}
                            render={({ field, fieldState }) => (
                                <Password
                                    id={field.name}
                                    {...field}
                                    toggleMask
                                    className={classNames({ "p-invalid": fieldState.error })}
                                    inputClassName="w-full"
                                    promptLabel="Şifre belirleyin"
                                    weakLabel="Zayıf" mediumLabel="Orta" strongLabel="Güçlü"
                                />
                            )}
                        />
                    </FormItem>
                    <FormItem
                        label="Yeni Şifre Tekrar"
                        name="confirmPassword"
                        required={isFieldRequired(changePasswordSchema, "confirmPassword")}
                    >
                        <Controller
                            name="confirmPassword"
                            control={methods.control}
                            render={({ field, fieldState }) => (
                                <Password
                                    id={field.name}
                                    {...field}
                                    feedback={false}
                                    toggleMask
                                    className={classNames({ "p-invalid": fieldState.error })}
                                    inputClassName="w-full"
                                />
                            )}
                        />
                    </FormItem>

                    <Button
                        type="submit" // onClick yerine type="submit"
                        label="Şifreyi Güncelle"
                        icon="pi pi-save"
                        loading={methods.formState.isSubmitting}
                        className="p-button-primary w-full mt-2"
                    />
                </form>
            </FormProvider >
        </Dialog >
    );
}