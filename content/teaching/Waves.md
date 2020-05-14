---
title: "Waves and Modern Physics"
date: 2016-12-17T20:04:57-07:00
draft: false
---

These are a set of course notes from a crash course for CEGEP Waves and
Modern Physics. They're missing sections on nuclear physics and are largely
just a summary of the important equations with very brief explanations in some
places. You can find a copy of the hand written notes at this
[link](/documents/WavesandModernPhysicsNotes.pdf).

I've also started to write out these notes below with a few more details where
necessary.

---

#  Part 1: Oscillations and Mechanical Waves

Any particle with a simple restoring force will result in simple harmonic
motion.

$$
F = -kx \Rightarrow x(t) = A cos\left(\omega t + \phi\right)
$$

where $\omega = \sqrt{k / m}$. How do you know? Well lets check. Using Newton's
second law, we can verify,

$$
F = ma = -kx \Rightarrow m \frac{d^2 x(t)}{dt^2} = -k x(t).
$$

Lets check our guess of simple harmonic motion ($x(t) = Acos(\omega t + \phi)$)
by replacing $x(t)$ on each side of our equation.

$$
\begin{align}
-m A \omega^2 \cos\left(\omega t + \phi\right) &= -k A \cos\left(\omega t + \phi\right) \\\\
A \omega^2 \cos\left(\omega t + \phi \right) &= \frac{k}{m} A \cos\left(\omega t + \phi \right) \\\\
A \omega^2 \cos\left(\omega t + \phi \right) &= A \omega^2 \cos\left(\omega t + \phi \right) \\,\\,\\,\\,\\, \checkmark
\end{align}
$$

This kind of force, and in turn this kind of motion, is extremely common. Some
examples are atoms in solids, diatomic molecules, buildings, and any elastic
motion. We call the constants $\omega$ and $\phi$ the angular velocity and the
phase respectively. The angular velocity is equivalent to two other metrics of
harmonic motion: the period and frequency.

The period, $T$, is the time required for a particle to return to its original
position,

$$
x(t) = x(t + T)
$$

If the particle returns to its original position, this means the argument of the
cosine has changed by $2\pi$ radians. This implies a relation to the angular
frequency,

$$
\omega (t + T) + \phi = \omega t + \phi + 2\pi \Rightarrow T = \frac{2\pi}{\omega}
$$

The frequency, $f$, is the number of cycles that happen per unit of time. For
instance, if the particle is in harmonic motion we might want to know how many
times it returns to its original position in 1 second. This is the frequency.
It is the inverse of the period,

$$
f = \frac{1}{T} = \frac{\omega}{2\pi}
$$

---

**Example - Car Driving over a Pot Hole**

Two people drive over a pothole in the road and the car starts to bounce up and
down afterwards because of the impact. If the car weighs 1300 kg, each person
 weighs 160 kg and the shocks have a spring tension of 20 kN / m, what is the
 frequency of the bouncing?

*Each spring exerts a force $-kx$ upwards on the car. In total this creates an
effective spring constant we use to calculate the frequency,*

$$
\begin{align}
F &= -\sum_{springs}kx = -\left(\sum_{springs} k\right) x = -k_{eff} x \\\\
k_{eff} &= 4 \times 20 kN / m = 80 \\,\rm{kN / m}, \\\\
m &= m_{car} + 2 \times m_{person} = 1460 \,\rm{kg}, \\\\
f &= \frac{1}{2\pi} \sqrt{\frac{k_{eff}}{m}} = 1.18 \\,\rm{Hz}.
\end{align}
$$

---

## Energy of a Harmonic Oscillator

The energy in a harmonic oscillator is,

$$
E = \frac{1}{2} k A ^2.
$$

How do we know? Well, there are a few ways to verify:

​	1.)	At the peak amplitude $x(t) = A$, there is no velocity so the energy is all potential,

$$
E = \frac{1}{2} kx^2 = \frac{1}{2} k A^2
$$

​	2.) At the x = 0 (equilibrium position) there is only kinetic energy:

$$
\begin{align}
E &= \frac{1}{2} m v^2 \\\\
   &= \frac{1}{2} m \left(\frac{d}{dt} \left(A \cos(\omega t + \phi\right)\right)^2 \\\\
   &= \frac{1}{2} m \left(- A \omega \sin(\omega t + \phi)\right)^2 \\\\
   &= \frac{1}{2} m A ^ 2 \omega ^ 2 \\\\
   &= \frac{1}{2} k A ^ 2
\end{align}
$$

Where the third line in the previous equation comes from the fact that
$x = 0 \Rightarrow \sin(\omega t + \phi) = 1$.

In general, we have,

$$
\begin{align}
E &= \frac{1}{2} m v^2 + \frac{1}{2} kx^2 \\\\
  &= \frac{1}{2} k A^2 (\sin^2(\omega t + \phi) + \cos^2(\omega t + \phi)) \\\\
  &= \frac{1}{2} k A^2
\end{align}
$$

While we've discussed systems of masses and springs here, there are other
examples of simple harmomic motion as well.

**The Pendulum**

In a pendulum a particle of mass $m$ swings on a cord of length $L$. The state is defined by the angled from vertical.

<center>
<img height="180px" alt="Pendulum" src="/img/waves/pendulum.svg"/>
</center>

For small angles the equation of motion is,

$$
\frac{d^2 \theta}{dt^2} = -\frac{g}{L} \theta.
$$

Comparing with the equation of motion for a particle on a spring we can see
they are the same with the following substitutions:

$$
\begin{align}
\theta(t) &= x(t) \\\\
k &= g \\\\
L &= m
\end{align}
$$

This means that simple harmonic motion is once again the solution with
$\omega = \sqrt{g / L}$.

## Damped Oscillations

If there is a damping force as well as a restoring force the motion becomes
"damped" harmonic motion.

$$
F = -kx - bv = ma
$$

[To be continued...]
