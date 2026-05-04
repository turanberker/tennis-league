import React, { useEffect, useRef, useState } from 'react';
import { Sidebar } from 'primereact/sidebar';
import { Controller, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { MatchScore, Status, UpdateScoreRequest } from '../../model/match.model';
import * as yup from 'yup';
import { getMatchInfo } from '../../api/matchService';
import { Toast } from 'primereact/toast';
import { Calendar } from 'primereact/calendar';
import { Button } from 'primereact/button';
import { InputNumber } from 'primereact/inputnumber';
import { Checkbox } from 'primereact/checkbox';
import Guard from '../../helper/Guard';

export interface MatchScoreSidebarProps {
    visible: boolean;
    matchId?: string;
    onHide: () => void;
    onSuccess: () => void;
    submitMatchScore: (matchId: string, score: MatchScore) => boolean | Promise<boolean>
}

// --- Validation Schemas ---
const baseScoreSchema = yup.object().shape({
    team1Score: yup.number().typeError('Sayı giriniz').required('Zorunlu').min(0).max(99),
    team2Score: yup.number().typeError('Sayı giriniz').required('Zorunlu').min(0).max(99),
}).test('tennis-set', 'İki tarafında skoru aynı olamaz', (value) => {
    if (!value) return false;
    const { team1Score, team2Score } = value;
    return (team1Score !== team2Score);
});

const superTieSchema = yup.object().shape({
    team1Score: yup.number().typeError('Sayı giriniz').required('Zorunlu').min(0).max(99),
    team2Score: yup.number().typeError('Sayı giriniz').required('Zorunlu').min(0).max(99),
}).test("super-tie", 'En az 2 fark olmalı', (value) => {
    if (!value) return true;
    const diff = Math.abs(value.team1Score - value.team2Score);
    return diff >= 2;
}).nullable().default(null);

export const matchScoreSchema = yup.object().shape({
    matchDate: yup.date().typeError('Geçerli bir tarih giriniz').required('Maç tarihi zorunludur'),
    set1: baseScoreSchema.required(),
    set2: baseScoreSchema.required(),
    superTie: superTieSchema,
}).test("super-tie-required", "Eşitlik durumunda süper tie skoru girmelisiniz", function (this: yup.TestContext, value) {
    debugger
    if (!value || !value.set1 || !value.set2) return true;
    const s1Winner = value.set1.team1Score > value.set1.team2Score

    const s2Winner = value.set2.team1Score > value.set2.team2Score

    // Durum 1-1 mi? (Setler paylaşıldı mı?)
    const isTie = s1Winner !== s2Winner;

    if (isTie) {
        // DURUM 1: Setler 1-1 ama Super Tie girilmemiş
        if (!value.superTie) {
            return this.createError({
                path: undefined,
                message: 'Setler berabere olduğu için Super Tie skoru girmelisiniz.'
            });
        }
    } else {
        // DURUM 2: Setler 2-0 veya 0-2 ama Super Tie alanı dolu/işaretli
        if (value.superTie !== null) {
            return this.createError({
                path: undefined,
                message: 'Maç 2-0 bittiği durumlarda Super Tie skoru girilemez.'
            });
        }
    }

    // Durum 2-0 veya 0-2 ise Super Tie zorunlu değil
    return true;
});

export function MatchScoreSidebar({ visible, onHide, matchId, onSuccess, submitMatchScore }: MatchScoreSidebarProps) {
    const [selectedMatch, setSelectedMatch] = useState<{ side1: string; side2: string } | null>(null);
    const [showSuperTie, setShowSuperTie] = useState(false);
    const calendarRef = useRef<Calendar>(null);
    const toast = useRef<Toast>(null);
    const [isReadOnly, setIsReadOnly] = useState(false);

    const { control, handleSubmit, setValue, formState: { errors, isSubmitting } } = useForm<UpdateScoreRequest>({
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
                const res = await getMatchInfo(matchId);
                if (res) {
                    setIsReadOnly(res.matchInfo.status === Status.SCORE_APPROVED)
                    setSelectedMatch({ side1: res.matchInfo.side1, side2: res.matchInfo.side2 });

                    setValue('matchDate', res.matchInfo.matchDate ? new Date(res.matchInfo.matchDate) : new Date());
                    setValue('set1', res.setScore.set1 || { team1Score: 0, team2Score: 0 });
                    setValue('set2', res.setScore.set2 || { team1Score: 0, team2Score: 0 });

                    if (res.setScore.superTie) {
                        setValue('superTie', res.setScore.superTie);
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
        console.log(data)
        const success = await submitMatchScore(matchId!, data);
        if (success) {
            toast.current?.show({ severity: 'success', summary: 'Başarılı', detail: 'Skor güncellendi' });
            onHide();
            onSuccess();
        }
    };

    return (
        <>
            <Toast ref={toast} />
            <Sidebar visible={visible} position="right" onHide={onHide} header={isReadOnly ? "Maç Skoru" : "Maç Skoru Düzenle"} className="p-sidebar-md">
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
                                    disabled={isReadOnly}
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
                    {[1, 2].map((setNum) => {
                        const setKey = `set${setNum}` as keyof UpdateScoreRequest;
                        const setError = errors[setKey];
                        return (
                            <div key={setNum} className="mb-4">
                                <h4 className="mb-2">Set {setNum}</h4>
                                <div className="grid">
                                    <div className="col-6">
                                        <label className="block mb-1 text-sm">{selectedMatch?.side1}</label>
                                        <Controller
                                            name={`set${setNum}.team1Score` as any}
                                            control={control}
                                            render={({ field }) => <InputNumber disabled={isReadOnly} value={field.value} onValueChange={(e) => field.onChange(e.value)} min={0} max={7} showButtons />}
                                        />
                                    </div>
                                    <div className="col-6">
                                        <label className="block mb-1 text-sm">{selectedMatch?.side2}</label>
                                        <Controller
                                            name={`set${setNum}.team2Score` as any}
                                            control={control}
                                            render={({ field }) => <InputNumber disabled={isReadOnly} value={field.value} onValueChange={(e) => field.onChange(e.value)} min={0} max={7} showButtons />}
                                        />
                                    </div>
                                </div>
                                {setError && (
                                    <small className="p-error block mt-1">
                                        <i className="pi pi-exclamation-triangle mr-2"></i>     {setError.message || (setError as any).root?.message}
                                    </small>
                                )}
                            </div>
                        )
                    }
                    )}

                    {/* SUPER TIE TOGGLE */}
                    <div className="field-checkbox my-4">
                        <Checkbox
                            inputId="stToggle"
                            disabled={isReadOnly}
                            checked={showSuperTie}
                            onChange={(e) => {
                                setShowSuperTie(e.checked || false);
                                setValue('superTie', e.checked ? { team1Score: 0, team2Score: 0 } : null);
                            }}
                        />
                        <label htmlFor="stToggle" className="ml-2">Super Tie Oynandı mı?</label>
                    </div>
                    {(errors.root?.message || (errors as any)[""]?.message) && (
                        <div className="p-error block mt-2 mb-2">
                            <small className="p-error flex align-items-center font-bold">
                                <i className="pi pi-exclamation-triangle mr-2"></i>
                                {errors.root?.message || (errors as any)[""]?.message}
                            </small>
                        </div>
                    )}

                    {/* SUPER TIE INPUTS */}
                    {showSuperTie && (
                        <div className="p-3 surface-ground border-round mb-4 " >
                            <h4 className="mt-0">Super Tie</h4>
                            <div className="grid">
                                <div className="col-6">
                                    <Controller
                                        name="superTie.team1Score"
                                        control={control}
                                        render={({ field }) => <InputNumber disabled={isReadOnly} value={field.value} className={errors.superTie ? 'p-invalid' : ''} onValueChange={(e) => field.onChange(e.value)} min={0} />}
                                    />
                                </div>
                                <div className="col-6">
                                    <Controller
                                        name="superTie.team2Score"
                                        control={control}
                                        render={({ field }) => <InputNumber disabled={isReadOnly} value={field.value} className={errors.superTie ? 'p-invalid' : ''} onValueChange={(e) => field.onChange(e.value)} min={0} />}
                                    />
                                </div>
                            </div>
                            {errors.superTie && (
                                <small className="p-error block mt-2">
                                    <i className="pi pi-exclamation-triangle mr-2"></i>    {(errors.superTie as any).message || (errors.superTie as any).root?.message}
                                </small>
                            )}
                        </div>


                    )}


                    <Guard onFail={() => setIsReadOnly(true)}>
                        {!isReadOnly && (
                            <Button type="submit" label="Değişiklikleri Kaydet" icon="pi pi-check" loading={isSubmitting} />
                        )}
                    </Guard>
                </form>
            </Sidebar>
        </>
    );
}