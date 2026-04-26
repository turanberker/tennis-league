import React, { useEffect, useRef, useState } from 'react';
import { Sidebar } from 'primereact/sidebar';
import { Controller, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { MatchScore } from '../../model/match.model';
import * as yup from 'yup';
import { getSetScores, updateMatchScore } from '../../api/matchService';
import { Toast } from 'primereact/toast';
import { Calendar } from 'primereact/calendar';
import { Button } from 'primereact/button';
import { InputNumber } from 'primereact/inputnumber';
import { Checkbox } from 'primereact/checkbox';

export interface MatchScoreSidebarProps {
    visible: boolean;
    matchId?: string;
    onHide: () => void;
    onSuccess: () => void;
}

// --- Validation Schemas ---
const baseScoreSchema = yup.object().shape({
    team1Score: yup.number().typeError('Sayı giriniz').required('Zorunlu').min(0).max(99),
    team2Score: yup.number().typeError('Sayı giriniz').required('Zorunlu').min(0).max(99),
});

const tennisSetSchema = baseScoreSchema.test('tennis-set', 'Geçerli set skoru giriniz', (value) => {
    if (!value) return false;
    const { team1Score, team2Score } = value;
    const max = Math.max(team1Score, team2Score);
    const min = Math.min(team1Score, team2Score);
    return (max === 6 && min <= 4) || (max === 7 && min === 5) || (max === 7 && min === 6);
});

const superTieSchema = yup.object({
    team1Score: yup.number().required('Zorunlu').min(0).max(99),
    team2Score: yup.number().required('Zorunlu').min(0).max(99),
}).test('super-tie', 'Min 10 ve 2 fark olmalı', (value) => {
    if (!value) return true;
    const max = Math.max(value.team1Score, value.team2Score);
    const diff = Math.abs(value.team1Score - value.team2Score);
    return max >= 10 && diff >= 2;
}).nullable().default(null);

export const matchScoreSchema = yup.object().shape({
    matchDate: yup.date().typeError('Geçerli bir tarih giriniz').required('Maç tarihi zorunludur'),
    set1: tennisSetSchema.required(),
    set2: tennisSetSchema.required(),
    superTie: superTieSchema,
});

export function MatchScoreSidebar({ visible, onHide, matchId, onSuccess }: MatchScoreSidebarProps) {
    const [selectedMatch, setSelectedMatch] = useState<{ side1: string; side2: string } | null>(null);
    const [showSuperTie, setShowSuperTie] = useState(false);
    const calendarRef = useRef<Calendar>(null);
    const toast = useRef<Toast>(null);

    const { control, handleSubmit, reset, setValue, formState: { errors, isSubmitting } } = useForm<MatchScore>({
        resolver: yupResolver(matchScoreSchema),
        defaultValues: {
            matchDate: null as any,
            set1: { team1Score: 0, team2Score: 0 },
            set2: { team1Score: 0, team2Score: 0 },
            superTie: null
        }
    });

    useEffect(() => {
        if (visible && matchId) {
            const loadData = async () => {
                const res = await getSetScores(matchId);
                if (res) {
                    setSelectedMatch({ side1: res.side1, side2: res.side2 });

                    // CRITICAL: String -> Date dönüşümü
                    setValue('matchDate', res.matchDate ? new Date(res.matchDate) : new Date());
                    setValue('set1', res.set1 || { team1Score: 0, team2Score: 0 });
                    setValue('set2', res.set2 || { team1Score: 0, team2Score: 0 });

                    if (res.superTie) {
                        setValue('superTie', res.superTie);
                        setShowSuperTie(true);
                    } else {
                        setValue('superTie', null);
                        setShowSuperTie(false);
                    }
                }
            };
            loadData();
        }
    }, [visible, matchId, setValue]);

    const onSubmit = async (data: MatchScore) => {
        const success = await updateMatchScore(matchId!, data);
        if (success) {
            toast.current?.show({ severity: 'success', summary: 'Başarılı', detail: 'Skor güncellendi' });
            onHide();
            onSuccess();
        }
    };

    return (
        <>
            <Toast ref={toast} />
            <Sidebar visible={visible} position="right" onHide={onHide} header="Maç Skoru Düzenle" className="p-sidebar-md">
                <form onSubmit={handleSubmit(onSubmit)} className="p-fluid">

                    {/* TARİH ALANI */}
                    <div className="field">
                        <label htmlFor="matchDate" className="font-bold">Maç Tarihi ve Saati</label>
                        <Controller
                            name="matchDate"
                            control={control}
                            render={({ field }) => (
                                <Calendar
                                    {...field}
                                    ref={calendarRef}
                                    id="matchDate"
                                    showTime
                                    hourFormat="24"
                                    stepMinute={10}

                                    className={errors.matchDate ? 'p-invalid' : ''}
                                    footerTemplate={() => (
                                        <div className="flex justify-content-end p-2">
                                            <Button type="button" label="Seç" className="p-button-sm" onClick={() => calendarRef.current?.hide()} />
                                        </div>
                                    )}
                                />
                            )}
                        />
                        {errors.matchDate && <small className="p-error">{errors.matchDate.message}</small>}
                    </div>

                    <hr className="my-4" />

                    {/* SETLER FONKSİYONU (Tekrarı önlemek için) */}
                    {[1, 2].map((setNum) => (
                        <div key={setNum} className="mb-4">
                            <h4 className="mb-2">Set {setNum}</h4>
                            <div className="grid">
                                <div className="col-6">
                                    <label className="block mb-1 text-sm">{selectedMatch?.side1}</label>
                                    <Controller
                                        name={`set${setNum}.team1Score` as any}
                                        control={control}
                                        render={({ field }) => <InputNumber value={field.value} onValueChange={(e) => field.onChange(e.value)} min={0} max={7} showButtons />}
                                    />
                                </div>
                                <div className="col-6">
                                    <label className="block mb-1 text-sm">{selectedMatch?.side2}</label>
                                    <Controller
                                        name={`set${setNum}.team2Score` as any}
                                        control={control}
                                        render={({ field }) => <InputNumber value={field.value} onValueChange={(e) => field.onChange(e.value)} min={0} max={7} showButtons />}
                                    />
                                </div>
                            </div>
                            {errors[`set${setNum}` as keyof MatchScore] && (
                                <small className="p-error">{(errors[`set${setNum}` as keyof MatchScore] as any)?.message}</small>
                            )}
                        </div>
                    ))}

                    {/* SUPER TIE TOGGLE */}
                    <div className="field-checkbox my-4">
                        <Checkbox
                            inputId="stToggle"
                            checked={showSuperTie}
                            onChange={(e) => {
                                setShowSuperTie(e.checked || false);
                                setValue('superTie', e.checked ? { team1Score: 0, team2Score: 0 } : null);
                            }}
                        />
                        <label htmlFor="stToggle" className="ml-2">Super Tie Oynandı mı?</label>
                    </div>

                    {/* SUPER TIE INPUTS */}
                    {showSuperTie && (
                        <div className="p-3 surface-ground border-round mb-4">
                            <h4 className="mt-0">Super Tie</h4>
                            <div className="grid">
                                <div className="col-6">
                                    <Controller
                                        name="superTie.team1Score"
                                        control={control}
                                        render={({ field }) => <InputNumber value={field.value} onValueChange={(e) => field.onChange(e.value)} min={0} />}
                                    />
                                </div>
                                <div className="col-6">
                                    <Controller
                                        name="superTie.team2Score"
                                        control={control}
                                        render={({ field }) => <InputNumber value={field.value} onValueChange={(e) => field.onChange(e.value)} min={0} />}
                                    />
                                </div>
                            </div>
                        </div>
                    )}

                    <Button type="submit" label="Değişiklikleri Kaydet" icon="pi pi-check" loading={isSubmitting} className="mt-2" />
                </form>
            </Sidebar>
        </>
    );
}