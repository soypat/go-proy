n = 100;
D = 4;
X = rand(D,n);
w = rand(D,1);
y_NoNoise = X.'*w;

mean = 1/(1+1/999)^999; %1/e
sigma_n = .222;
noise =mean+ sigma_n.*randn(n,1);

y = y_NoNoise+noise;
i_p =randi([1 n]);
i_q = randi([1 n]);
while i_p == i_q
    i_q = randi([1 n]);
end
xp = X(:,i_p);
xq = X(:,i_q);
COV = exp(-.5*(xp-xq).^2);
w_bar = (sigma_n^(-2))*( (sigma_n^(-2)) * X*X.' )^-1*X*y;

LinFunc = getLinFun(D);

f = X.'*w_bar;

plot (y)
hold on
plot (f)

function [linearFunction] = getLinFun(Dimension)
functionString = "@(w,x) ";

for i = 1:Dimension
    functionString = strcat(functionString,"w(",string(i),")*x(",string(i),")" );
    if i~=Dimension
        functionString = strcat(functionString,"+");
    end
end
linearFunction = eval(functionString);
end
