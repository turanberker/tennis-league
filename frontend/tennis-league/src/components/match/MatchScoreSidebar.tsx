import React, { useEffect, useRef, useState } from 'react';
import { Sidebar } from 'primereact/sidebar';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { MatchScore } from '../../model/match.model';
import * as yup from 'yup';
import { getSetScores, updateMatchScore } from '../../api/matchService';

import { Toast } from 'primereact/toast';

export interface MatchScoreSidebarProps {
    visible: boolean;
    matchId?: string;
    onHide: () => void;
    onSuccess: () => void;
}

/**
 * Ortak score alanı (gte=0,lte=99)
 */
const baseScoreSchema = yup.object().shape({
    team1Score: yup
        .number()
        .typeError('Sayı giriniz')
        .required('Zorunlu')
        .min(0, "0'dan küçük olamaz")
        .max(99, "99'dan büyük olamaz"),

    team2Score: yup
        .number()
        .typeError('Sayı giriniz')
        .required('Zorunlu')
        .min(0, "0'dan küçük olamaz")
        .max(99, "99'dan büyük olamaz"),
});

/**
 * tennis_set validation (backend: tennis_set)
 */
const tennisSetSchema = baseScoreSchema.test(
    'tennis-set',
    'Geçerli set skoru giriniz (6-0..4, 7-5, 7-6)',
    (value) => {
        if (!value) return false;

        const { team1Score, team2Score } = value;

        const max = Math.max(team1Score, team2Score);
        const min = Math.min(team1Score, team2Score);

        if (max === 6 && min <= 4) return true;
        if (max === 7 && min === 5) return true;
        if (max === 7 && min === 6) return true;

        return false;
    },
);
/**
 * super_tie validation (backend: super_tie)
 */
const superTieSchema = yup
    .object({
        team1Score: yup
            .number()
            .transform((v, o) => (o === '' ? undefined : v))
            .required('Zorunlu')
            .min(0)
            .max(99),

        team2Score: yup
            .number()
            .transform((v, o) => (o === '' ? undefined : v))
            .required('Zorunlu')
            .min(0)
            .max(99),
    })
    .test('super-tie', 'SuperTie min 10 ve 2 fark olmalı', (value) => {
        if (!value) return true;

        const max = Math.max(value.team1Score, value.team2Score);
        const diff = Math.abs(value.team1Score - value.team2Score);

        return max >= 10 && diff >= 2;
    })
    .nullable()
    .default(null);

export const matchScoreSchema = yup.object().shape({
    set1: tennisSetSchema.required(),
    set2: tennisSetSchema.required(),
    superTie: superTieSchema,
});

const initialValues = {
    set1: { team1Score: 0, team2Score: 0 },
    set2: { team1Score: 0, team2Score: 0 },
    superTie: null,
};

export function MatchScoreSidebar({ visible, onHide, matchId, onSuccess }: MatchScoreSidebarProps) {

    const [updateScoreVisible, setUpdateScoreVisible] = useState(false);
    const [selectedMatch, setSelectedMatch] = useState<{ side1: { name: string }, side2: { name: string } } | null>(null);
    const [showSuperTie, setShowSuperTie] = useState(false);

    useEffect(() => {
        const loadSetScores = async () => {
            if (matchId) {
                const setScores = await getSetScores(matchId);
                setSelectedMatch({
                    side1: { name: setScores.side1 },
                    side2: { name: setScores.side2 },
                });
                if (setScores) {
                    setValue('set1.team1Score', setScores.set1?.team1Score ?? null);
                    setValue('set1.team2Score', setScores.set1?.team2Score ?? null);

                    setValue('set2.team1Score', setScores.set2?.team1Score ?? null);
                    setValue('set2.team2Score', setScores.set2?.team2Score ?? null);
                    if (setScores.superTie) {
                        setValue('superTie.team1Score', setScores.superTie?.team1Score ?? null);
                        setValue('superTie.team2Score', setScores.superTie?.team2Score ?? null);
                        setShowSuperTie(true);
                    } else {
                        setShowSuperTie(false);
                    }
                }
            }

        }

        if (visible && matchId) {
            setUpdateScoreVisible(true);
            loadSetScores();
        } else {
            setUpdateScoreVisible(false);
        }
    }, [visible, matchId]);

    const toast = useRef<Toast>(null);
    const {
        register,
        handleSubmit,
        reset,
        setValue,
        formState: { errors, isSubmitting },
    } = useForm<MatchScore>({
        resolver: yupResolver(matchScoreSchema),
        defaultValues: initialValues,
    });

    const onHideHandler = () => {
        reset();
        setUpdateScoreVisible(false);
        onHide();
    };

    const onSubmit = async (data: MatchScore) => {

        const score = await updateMatchScore(matchId!, data);
        if (score) {
            console.log(score);
            toast.current?.show({
                severity: 'success',
                summary: 'Başarılı',
                detail: 'Maç Skoru Kaydedilmiştir',
                life: 3000,
            });

            reset();
            setUpdateScoreVisible(false);
            onSuccess();
        }

    };

    return (
        <>
            <Toast ref={toast} />

            <Sidebar
                header="Maç Skoru"
                visible={updateScoreVisible}
                position="right"
                onHide={() => { onHideHandler() }}
            >
                <form onSubmit={handleSubmit(onSubmit)} className="p-fluid">
                    {/* SET 1 */}
                    <h3>Set 1</h3>
                    <div className="p-grid">
                        <div className="p-col-6">
                            <label>{selectedMatch?.side1.name}</label>
                            <input
                                type="number"
                                max={7}
                                min={0}
                                {...register('set1.team1Score', { valueAsNumber: true })}
                                className="p-inputtext"
                            />
                        </div>

                        <div className="p-col-6">
                            <label>{selectedMatch?.side2.name}</label>
                            <input
                                type="number"
                                max={7}
                                min={0}
                                {...register('set1.team2Score', { valueAsNumber: true })}
                                className="p-inputtext"
                            />
                        </div>
                    </div>

                    {errors.set1?.message && (
                        <small className="p-error">{errors.set1.message}</small>
                    )}

                    {/* SET 2 */}
                    <h3 className="mt-4">Set 2</h3>
                    <div className="p-grid">
                        <div className="p-col-6">
                            <label>{selectedMatch?.side1.name}</label>
                            <input
                                type="number"
                                max={7}
                                min={0}
                                {...register('set2.team1Score', { valueAsNumber: true })}
                                className="p-inputtext"
                            />
                        </div>

                        <div className="p-col-6">
                            <label>{selectedMatch?.side2.name}</label>
                            <input
                                type="number"
                                max={7}
                                min={0}
                                {...register('set2.team2Score', { valueAsNumber: true })}
                                className="p-inputtext"
                            />
                        </div>
                    </div>
                    {errors.set2?.message && (
                        <small className="p-error">{errors.set2.message}</small>
                    )}

                    <div className="mt-4">
                        <div className="flex align-items-center gap-2">
                            <input
                                type="checkbox"
                                id="superTieToggle"
                                checked={showSuperTie}
                                onChange={(e) => {
                                    const checked = e.target.checked;
                                    setShowSuperTie(checked);

                                    if (checked) {
                                        setValue('superTie', {
                                            team1Score: 0,
                                            team2Score: 0,
                                        });
                                    } else {
                                        setValue('superTie', null);
                                    }
                                }}
                            />
                            <label htmlFor="superTieToggle">Super Tie Oynandı</label>
                        </div>
                    </div>

                    {/* SUPER TIE */}
                    {showSuperTie && (
                        <>
                            <h3 className="mt-3">Super Tie</h3>

                            <div className="p-grid">
                                <div className="p-col-6">
                                    <label>{selectedMatch?.side1.name}</label>
                                    <input
                                        type="number"
                                        {...register('superTie.team1Score', {
                                            valueAsNumber: true,
                                        })}
                                        className="p-inputtext"
                                    />
                                </div>

                                <div className="p-col-6">
                                    <label>{selectedMatch?.side2.name}</label>
                                    <input
                                        type="number"
                                        {...register('superTie.team2Score', {
                                            valueAsNumber: true,
                                        })}
                                        className="p-inputtext"
                                    />
                                </div>
                            </div>

                            {errors.superTie?.message && (
                                <small className="p-error">{errors.superTie.message}</small>
                            )}
                        </>
                    )}

                    {/* BUTTON */}
                    <div className="mt-4">
                        <button
                            type="submit"
                            className="p-button p-component"
                            disabled={isSubmitting}
                        >
                            Kaydet
                        </button>
                    </div>
                </form>
            </Sidebar>
        </>);
}